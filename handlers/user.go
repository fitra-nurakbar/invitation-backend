package handlers

import (
	"invitation-app/config"
	"invitation-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GET /users
func GetUsers(c *gin.Context) {
	var users []models.User
	result := config.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GET /users/:id
func GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var user models.User
	result := config.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// POST /users
func CreateUser(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:    uuid.New(),
		Email: input.Email,
		Name:  input.Name,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// PUT /users/:id
func UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var user models.User
	if result := config.DB.First(&user, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&user).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DELETE /users/:id
func DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result := config.DB.Delete(&models.User{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
