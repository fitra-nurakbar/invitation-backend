package models

import (
	"time"

	"github.com/google/uuid"
)

type UserTemplate struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TemplateID uuid.UUID `gorm:"type:uuid;not null" json:"template_id"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	GrantedAt  time.Time `gorm:"type:timestamptz;autoCreateTime" json:"granted_at"`

	// Relations
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template Template `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Order    Order    `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}
