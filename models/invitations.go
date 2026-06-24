package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/datatypes"
)

type InvitationStatus string

const (
    StatusActive  InvitationStatus = "active"
    StatusExpired InvitationStatus = "expired"
    StatusDraft   InvitationStatus = "draft"
)

type Invitation struct {
    ID         uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID     uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
    TemplateID uuid.UUID        `gorm:"type:uuid;not null" json:"template_id"`
    Slug       string           `gorm:"type:text;uniqueIndex;not null" json:"slug"`
    EventDate  time.Time        `gorm:"type:date" json:"event_date"`
    Status     InvitationStatus `gorm:"type:text;default:'draft'" json:"status"`
    ExpiresAt  time.Time        `gorm:"type:timestamptz" json:"expires_at"`
    Detail     datatypes.JSON   `gorm:"type:jsonb" json:"detail"`

    // Relations
    User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Template Template `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
    Messages []Message `gorm:"foreignKey:InvitationID" json:"messages,omitempty"`
}