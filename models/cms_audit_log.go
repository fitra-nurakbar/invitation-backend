package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type CMSAuditLog struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	InvitationID *uuid.UUID     `gorm:"type:uuid" json:"invitation_id,omitempty"`
	TemplateID   *uuid.UUID     `gorm:"type:uuid" json:"template_id,omitempty"`
	UserID       *uuid.UUID     `gorm:"type:uuid" json:"user_id,omitempty"`
	Action       string         `gorm:"type:text;not null" json:"action"`
	Metadata     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'::jsonb" json:"metadata"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
}
