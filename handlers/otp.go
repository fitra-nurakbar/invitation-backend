package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	"invitation-app/config"
	"invitation-app/models"
	"invitation-app/utils"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	otpPrefix       = "otp:"
	otpCooldownSecs = 60
	otpExpiryMins   = 10
	otpMaxAttempts  = 5
)

// POST /auth/otp/send — generate and send 6-digit OTP code to email
func SendOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email tidak valid"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Cooldown check (60 detik)
	cooldownKey := fmt.Sprintf("%scooldown:%s", otpPrefix, email)
	if exists := config.RedisClient.Exists(context.Background(), cooldownKey).Val(); exists > 0 {
		ttl := config.RedisClient.TTL(context.Background(), cooldownKey).Val()
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       "Mohon tunggu sebelum meminta kode baru",
			"retry_after": int(ttl.Seconds()),
		})
		return
	}

	// Cari atau auto-register user
	var user models.User
	result := config.DB.Where("LOWER(email) = ?", email).First(&user)
	if result.Error != nil {
		name := strings.TrimSpace(input.Name)
		if name == "" {
			for i, ch := range email {
				if ch == '@' {
					name = email[:i]
					break
				}
			}
		}

		user = models.User{
			ID:        uuid.New(),
			Email:     email,
			Name:      name,
			Password:  nil,
			Role:      models.RoleUser,
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun"})
			return
		}
		fmt.Printf("Auto-register user via OTP: %s\n", user.Email)
	}

	// Blokir admin dari login via OTP
	if user.IsAdmin() {
		c.JSON(http.StatusOK, gin.H{
			"message": "Jika email terdaftar, kode verifikasi akan dikirim ke email Anda",
		})
		return
	}

	// Generate 6-digit OTP
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kode verifikasi"})
		return
	}
	otpCode := fmt.Sprintf("%06d", n.Int64())

	// Simpan ke Redis
	otpKey := fmt.Sprintf("%scode:%s", otpPrefix, email)
	attemptsKey := fmt.Sprintf("%sattempts:%s", otpPrefix, email)

	if err := config.RedisClient.SetEx(
		context.Background(),
		otpKey,
		otpCode,
		time.Duration(otpExpiryMins)*time.Minute,
	).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kode verifikasi"})
		return
	}

	// Reset attempts
	config.RedisClient.Del(context.Background(), attemptsKey)

	// Set cooldown
	config.RedisClient.SetEx(
		context.Background(),
		cooldownKey,
		"1",
		time.Duration(otpCooldownSecs)*time.Second,
	)

	// Kirim email secara asinkron
	go func() {
		if err := emailSvc().SendOTP(user.Email, user.Name, otpCode); err != nil {
			fmt.Printf("Gagal kirim email OTP ke %s: %v\n", user.Email, err)
		} else {
			fmt.Printf("Kode OTP terkirim ke %s\n", user.Email)
		}
	}()

	response := gin.H{
		"message":    "Kode verifikasi telah dikirim ke email Anda",
		"expires_in": fmt.Sprintf("%d menit", otpExpiryMins),
		"cooldown_s": otpCooldownSecs,
	}

	// Di development, kirimkan otp_code langsung di response untuk mempermudah testing lokal
	if os.Getenv("APP_ENV") == "development" {
		response["otp_code"] = otpCode
	}

	c.JSON(http.StatusOK, response)
}

// POST /auth/otp/verify — verifikasi kode 6-digit OTP dan dapatkan JWT token
func VerifyOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan kode OTP 6-digit wajib diisi"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	otpCode := strings.TrimSpace(input.OTP)

	attemptsKey := fmt.Sprintf("%sattempts:%s", otpPrefix, email)
	otpKey := fmt.Sprintf("%scode:%s", otpPrefix, email)

	// Cek jumlah percobaan
	attempts, _ := config.RedisClient.Get(context.Background(), attemptsKey).Int()
	if attempts >= otpMaxAttempts {
		config.RedisClient.Del(context.Background(), otpKey)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Batas percobaan terlampaui. Silakan minta kode verifikasi baru.",
		})
		return
	}

	// Ambil kode tersimpan
	savedCode, err := config.RedisClient.Get(context.Background(), otpKey).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Kode verifikasi tidak valid atau sudah kedaluwarsa",
		})
		return
	}

	// Validasi kode
	if savedCode != otpCode {
		config.RedisClient.Incr(context.Background(), attemptsKey)
		config.RedisClient.Expire(context.Background(), attemptsKey, time.Duration(otpExpiryMins)*time.Minute)
		remaining := otpMaxAttempts - (attempts + 1)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("Kode verifikasi salah. Sisa percobaan: %d", remaining),
		})
		return
	}

	// Kode valid — bersihkan dari Redis
	config.RedisClient.Del(context.Background(), otpKey)
	config.RedisClient.Del(context.Background(), attemptsKey)

	// Ambil data user
	var user models.User
	if result := config.DB.Where("LOWER(email) = ?", email).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat sesi login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verifikasi berhasil",
		"token":   jwtToken,
		"user": gin.H{
			"id":     user.ID,
			"email":  user.Email,
			"name":   user.Name,
			"role":   user.Role,
			"avatar": user.Avatar,
		},
	})
}
