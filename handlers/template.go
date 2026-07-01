package handlers

import (
	"invitation-app/config"
	"invitation-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GET /templates/public
func GetPublicTemplates(c *gin.Context) {
	var templates []models.Template
	config.DB.Where("is_active = ?", true).
		Order("price asc").
		Find(&templates)
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

// GET /templates
func GetTemplates(c *gin.Context) {
	var templates []models.Template
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

// GET /templates/:id
func GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var template models.Template
	if result := config.DB.First(&template, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

// POST /templates
func CreateTemplate(c *gin.Context) {
	var input struct {
		Name              string `json:"name" binding:"required"`
		Price             int    `json:"price" binding:"required"`
		IsActive          bool   `json:"is_active"`
		OrderDeadlineDays int    `json:"order_deadline_days"`
		ActiveDaysAfter   int    `json:"active_days_after"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := models.Template{
		ID:                uuid.New(),
		Name:              input.Name,
		Price:             input.Price,
		IsActive:          input.IsActive,
		OrderDeadlineDays: input.OrderDeadlineDays,
		ActiveDaysAfter:   input.ActiveDaysAfter,
	}

	config.DB.Create(&template)
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

// PUT /templates/:id
func UpdateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var template models.Template
	if result := config.DB.First(&template, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template tidak ditemukan"})
		return
	}

	var input struct {
		Name              string `json:"name"`
		Price             int    `json:"price"`
		IsActive          bool   `json:"is_active"`
		OrderDeadlineDays int    `json:"order_deadline_days"`
		ActiveDaysAfter   int    `json:"active_days_after"`
	}
	c.ShouldBindJSON(&input)
	config.DB.Model(&template).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": template})
}

// DELETE /templates/:id
func DeleteTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Soft delete — set is_active = false
	result := config.DB.Model(&models.Template{}).
		Where("id = ?", id).
		Update("is_active", false)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template dinonaktifkan"})
}
