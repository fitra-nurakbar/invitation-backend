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
		Name    string `json:"name" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil IP address tamu
	ipAddress := c.ClientIP()

	message := models.Message{
		ID:           uuid.New(),
		InvitationID: invitation.ID,
		Name:         input.Name,
		Message:      input.Message,
		IPAddress:    ipAddress,
	}

	config.DB.Create(&message)
	c.JSON(http.StatusCreated, gin.H{"data": message})
}

func GetAllMessages(c *gin.Context) {
	var messages []models.Message
	result := config.DB.Find(&messages)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// DELETE /messages/:id
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
