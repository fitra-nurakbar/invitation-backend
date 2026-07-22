package handlers

import (
	"context"
	"fmt"
	"invitation-app/config"
	"invitation-app/models"
	"invitation-app/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const magicLinkPrefix = "magic_link:"

// POST /auth/magic-link — request magic link (auto-register jika belum ada)
func RequestMagicLink(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name"` // opsional — untuk auto-register
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek cooldown
	cooldownKey := fmt.Sprintf("%scooldown:%s", magicLinkPrefix, input.Email)
	if exists := config.RedisClient.Exists(context.Background(), cooldownKey).Val(); exists > 0 {
		ttl := config.RedisClient.TTL(context.Background(), cooldownKey).Val()
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       "Mohon tunggu sebelum request link baru",
			"retry_after": int(ttl.Seconds()),
		})
		return
	}

	// Cari atau buat user (auto-register)
	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		// User belum ada → auto-register sebagai user biasa
		name := input.Name
		if name == "" {
			// Gunakan bagian sebelum @ sebagai nama default
			for i, c := range input.Email {
				if c == '@' {
					name = input.Email[:i]
					break
				}
			}
		}

		user = models.User{
			ID:        uuid.New(),
			Email:     input.Email,
			Name:      name,
			Password:  nil, // user biasa tidak punya password
			Role:      models.RoleUser,
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun"})
			return
		}

		fmt.Printf("✅ Auto-register user baru: %s\n", user.Email)
	}

	// Blokir admin dari magic link
	if user.IsAdmin() {
		// Tetap return pesan generik — jangan bocorkan info
		c.JSON(http.StatusOK, gin.H{
			"message": "Jika email terdaftar, link login akan dikirim ke email kamu",
		})
		return
	}

	// Generate token
	token := uuid.New().String()

	expiredMinutes, _ := strconv.Atoi(os.Getenv("MAGIC_LINK_EXPIRED_MINUTES"))
	if expiredMinutes == 0 {
		expiredMinutes = 15
	}

	// Simpan token ke Redis
	tokenKey := fmt.Sprintf("%stoken:%s", magicLinkPrefix, token)
	if err := config.RedisClient.SetEx(
		context.Background(),
		tokenKey,
		user.ID.String(),
		time.Duration(expiredMinutes)*time.Minute,
	).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate magic link"})
		return
	}

	// Set cooldown 1 menit
	config.RedisClient.SetEx(
		context.Background(),
		cooldownKey,
		"1",
		1*time.Minute,
	)

	// Buat magic URL
	baseURL := os.Getenv("MAGIC_LINK_BASE_URL")
	magicURL := fmt.Sprintf("%s/auth/verify?token=%s", baseURL, token)

	// Kirim email
	go func() {
		if err := emailSvc().SendMagicLink(user.Email, user.Name, magicURL); err != nil {
			fmt.Println("❌ Gagal kirim magic link:", err)
		} else {
			fmt.Printf("✅ Magic link terkirim ke %s\n", user.Email)
		}
	}()

	response := gin.H{
		"message": "Jika email terdaftar, link login akan dikirim ke email kamu",
	}
	if os.Getenv("APP_ENV") == "development" {
		response["magic_url"] = magicURL
		response["expires_in"] = fmt.Sprintf("%d menit", expiredMinutes)
	}

	c.JSON(http.StatusOK, response)
}

// GET /auth/verify?token=xxx — verify → redirect ke frontend
func VerifyMagicLink(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token tidak ditemukan"})
		return
	}

	tokenKey := fmt.Sprintf("%stoken:%s", magicLinkPrefix, token)
	userIDStr, err := config.RedisClient.Get(context.Background(), tokenKey).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Link tidak valid atau sudah expired"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Hapus token — one time use
	config.RedisClient.Del(context.Background(), tokenKey)

	var user models.User
	if result := config.DB.First(&user, "id = ?", userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	frontendURL := os.Getenv("MAGIC_LINK_BASE_URL")
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, jwtToken)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GET /auth/verify-api?token=xxx — verify → JSON (untuk SPA/mobile)
func VerifyMagicLinkAPI(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token tidak ditemukan"})
		return
	}

	tokenKey := fmt.Sprintf("%stoken:%s", magicLinkPrefix, token)
	userIDStr, err := config.RedisClient.Get(context.Background(), tokenKey).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Link tidak valid atau sudah expired"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Hapus token — one time use
	config.RedisClient.Del(context.Background(), tokenKey)

	var user models.User
	if result := config.DB.First(&user, "id = ?", userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
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
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}
