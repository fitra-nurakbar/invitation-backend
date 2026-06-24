package handlers

import (
	"encoding/json"
	"invitation-app/config"
	"invitation-app/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// GET /invitations
func GetInvitations(c *gin.Context) {
	var invitations []models.Invitation
	config.DB.Preload("User").Preload("Template").Find(&invitations)
	c.JSON(http.StatusOK, gin.H{"data": invitations})
}

// GET /invitations/:slug
func GetInvitationBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var invitation models.Invitation
	result := config.DB.
		Preload("User").
		Preload("Template").
		Preload("Messages").
		Where("slug = ?", slug).
		First(&invitation)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": invitation})
}

// POST /invitations
func CreateInvitation(c *gin.Context) {
	var input struct {
		UserID     string                 `json:"user_id" binding:"required"`
		TemplateID string                 `json:"template_id" binding:"required"`
		Slug       string                 `json:"slug" binding:"required"`
		EventDate  string                 `json:"event_date" binding:"required"` // format: 2006-01-02
		Detail     map[string]interface{} `json:"detail"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(input.UserID)
	templateID, _ := uuid.Parse(input.TemplateID)
	eventDate, _ := time.Parse("2006-01-02", input.EventDate)

	// Ambil template untuk hitung expires_at
	var template models.Template
	config.DB.First(&template, "id = ?", templateID)
	expiresAt := eventDate.AddDate(0, 0, template.ActiveDaysAfter)

	// Convert detail map ke jsonb
	detailJSON, _ := json.Marshal(input.Detail)

	invitation := models.Invitation{
		ID:         uuid.New(),
		UserID:     userID,
		TemplateID: templateID,
		Slug:       input.Slug,
		EventDate:  eventDate,
		Status:     models.StatusDraft,
		ExpiresAt:  expiresAt,
		Detail:     datatypes.JSON(detailJSON),
	}

	result := config.DB.Create(&invitation)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": invitation})
}

// PUT /invitations/:id
func UpdateInvitation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var invitation models.Invitation
	if result := config.DB.First(&invitation, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	var input struct {
		Status    string                 `json:"status"`
		EventDate string                 `json:"event_date"`
		Detail    map[string]interface{} `json:"detail"`
	}
	c.ShouldBindJSON(&input)

	updates := map[string]interface{}{}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.EventDate != "" {
		eventDate, _ := time.Parse("2006-01-02", input.EventDate)
		updates["event_date"] = eventDate
	}
	if input.Detail != nil {
		detailJSON, _ := json.Marshal(input.Detail)
		updates["detail"] = datatypes.JSON(detailJSON)
	}

	config.DB.Model(&invitation).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"data": invitation})
}

// DELETE /invitations/:id
func DeleteInvitation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result := config.DB.Delete(&models.Invitation{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Undangan berhasil dihapus"})
}
