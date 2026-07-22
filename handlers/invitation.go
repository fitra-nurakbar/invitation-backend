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

// GET /admin/invitations
func GetInvitations(c *gin.Context) {
	var invitations []models.Invitation
	config.DB.
		Preload("User").
		Preload("Template").
		Order("created_at DESC").
		Find(&invitations)
	c.JSON(http.StatusOK, gin.H{"data": invitations})
}

// GET /invitations/my
func GetMyInvitations(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var invitations []models.Invitation
	config.DB.
		Preload("Template").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&invitations)

	c.JSON(http.StatusOK, gin.H{"data": invitations})
}

// GET /invitations/my/:id
func GetMyInvitation(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var invitation models.Invitation
	result := config.DB.
		Preload("Template").
		Preload("Messages").
		Where("id = ? AND user_id = ?", id, userID).
		First(&invitation)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": invitation})
}

// GET /invitations/:slug — publik
func GetInvitationBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var invitation models.Invitation
	result := config.DB.
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
	userID := c.MustGet("user_id").(uuid.UUID)

	var input struct {
		TemplateID string                     `json:"template_id" binding:"required"`
		Slug       string                     `json:"slug" binding:"required"`
		EventDate  string                     `json:"event_date" binding:"required"`
		Detail     models.InvitationDetailInput `json:"detail"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templateID, err := uuid.Parse(input.TemplateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template_id tidak valid"})
		return
	}

	eventDate, err := time.Parse("2006-01-02", input.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format event_date tidak valid, gunakan YYYY-MM-DD"})
		return
	}

	// Validasi field wajib
	if input.Detail.Groom.Name == "" || input.Detail.Bride.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama mempelai pria dan wanita wajib diisi"})
		return
	}

	if len(input.Detail.WeddingEvent) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimal satu acara pernikahan wajib diisi"})
		return
	}

	// Cek slug belum dipakai
	var existing models.Invitation
	if result := config.DB.Where("slug = ?", input.Slug).First(&existing); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Slug sudah digunakan, pilih slug lain"})
		return
	}

	// Cek user punya akses template
	var userTemplate models.UserTemplate
	if result := config.DB.Where(
		"user_id = ? AND template_id = ?",
		userID, templateID,
	).First(&userTemplate); result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Kamu belum memiliki template ini, silakan beli terlebih dahulu",
		})
		return
	}

	// Ambil template untuk hitung expires_at
	var template models.Template
	if result := config.DB.First(&template, "id = ? AND is_active = true", templateID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template tidak ditemukan atau tidak aktif"})
		return
	}

	expiresAt := eventDate.AddDate(0, 0, template.ActiveDaysAfter)

	// Convert detail ke JSON
	detailJSON, err := json.Marshal(input.Detail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses detail undangan"})
		return
	}

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

	if err := config.DB.Create(&invitation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat undangan"})
		return
	}

	// Preload template untuk response
	config.DB.Preload("Template").First(&invitation, "id = ?", invitation.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Undangan berhasil dibuat",
		"data":    invitation,
	})
}

// DELETE /invitations/delete/:id
func DeleteInvitation(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Hapus messages dulu
	config.DB.Where("invitation_id = ?", id).Delete(&models.Message{})

	// Hapus invitation milik user
	result := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Invitation{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Undangan berhasil dihapus"})
}