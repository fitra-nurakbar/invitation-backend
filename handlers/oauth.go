package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"invitation-app/config"
	"invitation-app/models"
	"invitation-app/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// GET /auth/google — redirect ke halaman login Google
func GoogleLogin(c *gin.Context) {
	// State untuk CSRF protection — simpan di Redis
	state := uuid.New().String()

	stateKey := fmt.Sprintf("oauth_state:%s", state)
	config.RedisClient.SetEx(
		context.Background(),
		stateKey,
		"valid",
		10*time.Minute, // state expire 10 menit
	)

	url := config.GoogleOAuthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /auth/google/callback — callback dari Google
func GoogleCallback(c *gin.Context) {
	// Verifikasi state — anti CSRF
	state := c.Query("state")
	stateKey := fmt.Sprintf("oauth_state:%s", state)

	exists := config.RedisClient.Exists(context.Background(), stateKey).Val()
	if exists == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "State tidak valid atau expired"})
		return
	}

	// Hapus state setelah dipakai
	config.RedisClient.Del(context.Background(), stateKey)

	// Ambil code dari query
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code tidak ditemukan"})
		return
	}

	// Tukar code dengan token Google
	googleToken, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal exchange token Google"})
		return
	}

	// Ambil data user dari Google
	googleUser, err := getGoogleUserInfo(googleToken.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data user Google"})
		return
	}

	// Pastikan email sudah terverifikasi Google
	if !googleUser.VerifiedEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email Google belum terverifikasi"})
		return
	}

	// Cari atau buat user
	user, isNew, err := findOrCreateGoogleUser(googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses akun"})
		return
	}

	if isNew {
		fmt.Printf("✅ User baru via Google: %s\n", user.Email)
	} else {
		fmt.Printf("✅ User login via Google: %s\n", user.Email)
	}

	// Generate JWT
	jwtToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	// Redirect ke frontend dengan JWT
	frontendURL := os.Getenv("MAGIC_LINK_BASE_URL")
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, jwtToken)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GoogleCallbackAPI — untuk SPA yang tidak mau redirect
func GoogleCallbackAPI(c *gin.Context) {
	state := c.Query("state")
	stateKey := fmt.Sprintf("oauth_state:%s", state)

	exists := config.RedisClient.Exists(context.Background(), stateKey).Val()
	if exists == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "State tidak valid atau expired"})
		return
	}
	config.RedisClient.Del(context.Background(), stateKey)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code tidak ditemukan"})
		return
	}

	googleToken, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal exchange token Google"})
		return
	}

	googleUser, err := getGoogleUserInfo(googleToken.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data user Google"})
		return
	}

	if !googleUser.VerifiedEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email Google belum terverifikasi"})
		return
	}

	user, _, err := findOrCreateGoogleUser(googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses akun"})
		return
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   jwtToken,
		"user": gin.H{
			"id":     user.ID,
			"email":  user.Email,
			"name":   user.Name,
			"avatar": user.Avatar,
			"role":   user.Role,
		},
	})
}

// ─── HELPER ──────────────────────────────────────────────────

// getGoogleUserInfo — ambil data user dari Google API
func getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("gagal request Google userinfo: %w", err)
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("gagal decode Google userinfo: %w", err)
	}

	return &userInfo, nil
}

// findOrCreateGoogleUser — cari user by google_id atau email, buat jika belum ada
func findOrCreateGoogleUser(googleUser *GoogleUserInfo) (*models.User, bool, error) {
	var user models.User
	isNew := false

	// 1. Cari by google_id dulu (returning user via Google)
	result := config.DB.Where("google_id = ?", googleUser.ID).First(&user)
	if result.Error == nil {
		// Update avatar jika berubah
		if user.Avatar == nil || (user.Avatar != nil && *user.Avatar != googleUser.Picture) {
			config.DB.Model(&user).Update("avatar", googleUser.Picture)
			user.Avatar = &googleUser.Picture
		}
		return &user, false, nil
	}

	// 2. Cari by email (user sudah ada via magic link)
	result = config.DB.Where("email = ?", googleUser.Email).First(&user)
	if result.Error == nil {
		// Link google_id ke akun yang sudah ada
		config.DB.Model(&user).Updates(map[string]interface{}{
			"google_id": googleUser.ID,
			"avatar":    googleUser.Picture,
		})
		user.GoogleID = &googleUser.ID
		user.Avatar = &googleUser.Picture
		return &user, false, nil
	}

	// 3. Buat user baru (pertama kali login via Google)
	googleID := googleUser.ID
	avatar := googleUser.Picture

	user = models.User{
		ID:        uuid.New(),
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		GoogleID:  &googleID,
		Avatar:    &avatar,
		Password:  nil,
		Role:      models.RoleUser,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, false, fmt.Errorf("gagal buat user baru: %w", err)
	}

	isNew = true
	return &user, isNew, nil
}