package handlers

import (
	"invitation-app/config"
	"invitation-app/models"
	"invitation-app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// POST /admin/auth/login — khusus admin
func AdminLogin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user by email
	var user models.User
	if result := config.DB.Where("email = ? AND role = ?", input.Email, models.RoleAdmin).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// Pastikan punya password
	if !user.HasPassword() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun ini tidak menggunakan password"})
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// POST /admin/auth/create — buat akun admin baru (hanya oleh admin)
func CreateAdmin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek email sudah ada
	var existing models.User
	if result := config.DB.Where("email = ?", input.Email).First(&existing); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	pwd := string(hashedPassword)
	admin := models.User{
		ID:        uuid.New(),
		Email:     input.Email,
		Name:      input.Name,
		Password:  &pwd,
		Role:      models.RoleAdmin,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Akun admin berhasil dibuat",
		"data": gin.H{
			"id":    admin.ID,
			"email": admin.Email,
			"name":  admin.Name,
			"role":  admin.Role,
		},
	})
}

// GET /auth/me — profile user yang sedang login
func Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var user models.User
	if result := config.DB.First(&user, "id = ?", userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}

// POST /auth/change-password — khusus admin
func ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var user models.User
	config.DB.First(&user, "id = ?", userID)

	// Pastikan admin
	if !user.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa ganti password"})
		return
	}

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !user.HasPassword() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Akun ini tidak menggunakan password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama tidak sesuai"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	pwd := string(hashed)
	config.DB.Model(&user).Update("password", &pwd)

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}