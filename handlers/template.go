package handlers

import (
	"encoding/json"
	"invitation-app/config"
	"invitation-app/models"
	"net/http"
	"strings"

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
	result := config.DB.Find(&templates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
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

func validateTemplateSchema(themeRaw, blockRaw json.RawMessage) (models.ThemeConfig, models.BlockConfig, error) {
	theme := models.DefaultClassicThemeConfig()
	blocks := models.DefaultStandardBlockConfig()
	if len(themeRaw) > 0 && string(themeRaw) != "null" {
		if err := json.Unmarshal(themeRaw, &theme); err != nil {
			return theme, blocks, err
		}
	}
	if len(blockRaw) > 0 && string(blockRaw) != "null" {
		if err := json.Unmarshal(blockRaw, &blocks); err != nil {
			return theme, blocks, err
		}
	}
	if err := models.ValidateThemeConfig(&theme); err != nil {
		return theme, blocks, err
	}
	if err := models.ValidateBlockConfig(&blocks); err != nil {
		return theme, blocks, err
	}
	return theme, blocks, nil
}

// POST /templates
func CreateTemplate(c *gin.Context) {
	var input struct {
		Name              string          `json:"name" binding:"required"`
		Slug              string          `json:"slug" binding:"required"`
		Description       string          `json:"description"`
		Thumbnail         string          `json:"thumbnail"`
		PreviewURL        string          `json:"preview_url"`
		Category          string          `json:"category"`
		Price             int             `json:"price"`
		IsPremium         bool            `json:"is_premium"`
		IsActive          bool            `json:"is_active"`
		OrderDeadlineDays int             `json:"order_deadline_days"`
		ActiveDaysAfter   int             `json:"active_days_after"`
		ThemeConfig       json.RawMessage `json:"theme_config"`
		BlockConfig       json.RawMessage `json:"block_config"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Category == "" {
		input.Category = "wedding"
	}
	theme, blocks, err := validateTemplateSchema(input.ThemeConfig, input.BlockConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "schema template tidak valid: " + err.Error()})
		return
	}

	template := models.Template{
		ID: uuid.New(), Slug: strings.TrimSpace(input.Slug), Name: strings.TrimSpace(input.Name),
		Description: input.Description, Thumbnail: input.Thumbnail, PreviewURL: input.PreviewURL,
		Category: input.Category, Price: input.Price, IsPremium: input.IsPremium, IsActive: input.IsActive,
		OrderDeadlineDays: input.OrderDeadlineDays, ActiveDaysAfter: input.ActiveDaysAfter,
		ThemeConfig: models.ToJSONBytes(theme), BlockConfig: models.ToJSONBytes(blocks), Version: 1,
	}
	if template.Slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug wajib diisi"})
		return
	}
	if result := config.DB.Create(&template); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

// POST /templates (legacy implementation replaced above)
func createTemplateLegacyRemoved(c *gin.Context) {
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
		Name              *string         `json:"name"`
		Slug              *string         `json:"slug"`
		Description       *string         `json:"description"`
		Thumbnail         *string         `json:"thumbnail"`
		PreviewURL        *string         `json:"preview_url"`
		Category          *string         `json:"category"`
		Price             *int            `json:"price"`
		IsPremium         *bool           `json:"is_premium"`
		IsActive          *bool           `json:"is_active"`
		OrderDeadlineDays *int            `json:"order_deadline_days"`
		ActiveDaysAfter   *int            `json:"active_days_after"`
		ThemeConfig       json.RawMessage `json:"theme_config"`
		BlockConfig       json.RawMessage `json:"block_config"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Slug != nil && strings.TrimSpace(*input.Slug) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug tidak boleh kosong"})
		return
	}
	if len(input.ThemeConfig) > 0 || len(input.BlockConfig) > 0 {
		themeRaw, blockRaw := input.ThemeConfig, input.BlockConfig
		if len(themeRaw) == 0 {
			themeRaw = json.RawMessage(template.ThemeConfig)
		}
		if len(blockRaw) == 0 {
			blockRaw = json.RawMessage(template.BlockConfig)
		}
		theme, blocks, schemaErr := validateTemplateSchema(themeRaw, blockRaw)
		if schemaErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "schema template tidak valid: " + schemaErr.Error()})
			return
		}
		template.ThemeConfig = models.ToJSONBytes(theme)
		template.BlockConfig = models.ToJSONBytes(blocks)
		template.Version++
	}
	if input.Name != nil {
		template.Name = strings.TrimSpace(*input.Name)
	}
	if input.Slug != nil {
		template.Slug = strings.TrimSpace(*input.Slug)
	}
	if input.Description != nil {
		template.Description = *input.Description
	}
	if input.Thumbnail != nil {
		template.Thumbnail = *input.Thumbnail
	}
	if input.PreviewURL != nil {
		template.PreviewURL = *input.PreviewURL
	}
	if input.Category != nil {
		template.Category = *input.Category
	}
	if input.Price != nil {
		template.Price = *input.Price
	}
	if input.IsPremium != nil {
		template.IsPremium = *input.IsPremium
	}
	if input.IsActive != nil {
		template.IsActive = *input.IsActive
	}
	if input.OrderDeadlineDays != nil {
		template.OrderDeadlineDays = *input.OrderDeadlineDays
	}
	if input.ActiveDaysAfter != nil {
		template.ActiveDaysAfter = *input.ActiveDaysAfter
	}
	if result := config.DB.Save(&template); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
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
