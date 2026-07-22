package services

import (
	"fmt"
	"invitation-app/config"
	"invitation-app/models"

	"github.com/google/uuid"
)

type TemplateAccessService struct{}

func NewTemplateAccessService() *TemplateAccessService {
	return &TemplateAccessService{}
}

// Grant — berikan akses template ke user setelah bayar
func (s *TemplateAccessService) Grant(userID, templateID, orderID uuid.UUID) error {
	// Cek apakah sudah punya akses
	var existing models.UserTemplate
	result := config.DB.Where(
		"user_id = ? AND template_id = ?",
		userID, templateID,
	).First(&existing)

	if result.Error == nil {
		// Sudah punya akses — idempotent, tidak error
		return nil
	}

	access := models.UserTemplate{
		ID:         uuid.New(),
		UserID:     userID,
		TemplateID: templateID,
		OrderID:    orderID,
	}

	if err := config.DB.Create(&access).Error; err != nil {
		return fmt.Errorf("gagal grant akses template: %w", err)
	}

	return nil
}

// HasAccess — cek apakah user punya akses ke template
func (s *TemplateAccessService) HasAccess(userID, templateID uuid.UUID) bool {
	var access models.UserTemplate
	result := config.DB.Where(
		"user_id = ? AND template_id = ?",
		userID, templateID,
	).First(&access)
	return result.Error == nil
}

// GetUserTemplates — ambil semua template yang dimiliki user
func (s *TemplateAccessService) GetUserTemplates(userID uuid.UUID) ([]models.UserTemplate, error) {
	var accesses []models.UserTemplate
	result := config.DB.
		Preload("Template").
		Where("user_id = ?", userID).
		Order("granted_at DESC").
		Find(&accesses)

	if result.Error != nil {
		return nil, fmt.Errorf("gagal ambil template user: %w", result.Error)
	}
	return accesses, nil
}
