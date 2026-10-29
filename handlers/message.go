package handlers

import (
	"invitation-app/config"
	"invitation-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GET /invitations/:slug/messages
func GetMessages(c *gin.Context) {
	slug := c.Param("slug")

	var invitation models.Invitation
	if result := config.DB.Where("slug = ?", slug).First(&invitation); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	var messages []models.Message
	config.DB.Where("invitation_id = ?", invitation.ID).
		Order("created_at DESC").
		Find(&messages)

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// POST /invitations/:slug/messages
func CreateMessage(c *gin.Context) {
	slug := c.Param("slug")

	var invitation models.Invitation
	if result := config.DB.Where("slug = ?", slug).First(&invitation); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	var input struct {
		Name       string `json:"name" binding:"required"`
		Message    string `json:"message" binding:"required"`
		Attendance string `json:"attendance"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendance := input.Attendance
	if attendance != "tidak_hadir" {
		attendance = "hadir"
	}

	// Ambil IP address tamu
	ipAddress := c.ClientIP()

	message := models.Message{
		ID:           uuid.New(),
		InvitationID: invitation.ID,
		Name:         input.Name,
		Message:      input.Message,
		Attendance:   attendance,
		IPAddress:    ipAddress,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": message})
}

// GET /invitations/me/:id/messages
func GetMyInvitationMessages(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	invitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID undangan tidak valid"})
		return
	}

	var invitation models.Invitation
	if err := config.DB.Where("id = ? AND user_id = ?", invitationID, userID).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	var messages []models.Message
	config.DB.Where("invitation_id = ?", invitation.ID).Order("created_at DESC").Find(&messages)

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// DELETE /invitations/me/:id/messages/:msg_id
func DeleteMyInvitationMessage(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	invitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID undangan tidak valid"})
		return
	}

	msgID, err := uuid.Parse(c.Param("msg_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesan tidak valid"})
		return
	}

	var invitation models.Invitation
	if err := config.DB.Where("id = ? AND user_id = ?", invitationID, userID).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	result := config.DB.Where("id = ? AND invitation_id = ?", msgID, invitation.ID).Delete(&models.Message{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pesan ucapan berhasil dihapus"})
}

// GET /admin/messages
func GetAllMessages(c *gin.Context) {
	var messages []models.Message
	result := config.DB.Preload("Invitation").Order("created_at DESC").Find(&messages)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// DELETE /admin/messages/:id
func DeleteMessage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result := config.DB.Delete(&models.Message{}, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pesan berhasil dihapus"})
}
