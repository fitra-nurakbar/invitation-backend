package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	InvitationID uuid.UUID `gorm:"type:uuid;not null" json:"invitation_id"`
	Name         string    `gorm:"type:text" json:"name"`
	Message      string    `gorm:"type:text" json:"message"`
	IPAddress    string    `gorm:"type:inet" json:"ip_address"`
	CreatedAt    time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relation
	Invitation Invitation `gorm:"foreignKey:InvitationID" json:"invitation,omitempty"`
}

// Custom scanner untuk type inet dari PostgreSQL
func (m *Message) AfterFind(tx *gorm.DB) error {
	return nil
}
