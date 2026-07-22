package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusPaid    OrderStatus = "paid"
	OrderStatusExpired OrderStatus = "expired"
	OrderStatusFailed  OrderStatus = "failed"
)

type Order struct {
	ID             uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	TemplateID     uuid.UUID   `gorm:"type:uuid;not null" json:"template_id"`
	InvoiceID      string      `gorm:"type:text;uniqueIndex;not null" json:"invoice_id"`
	InvoiceURL     string      `gorm:"type:text;not null" json:"invoice_url"`
	Amount         int         `gorm:"type:integer;not null" json:"amount"`
	Status         OrderStatus `gorm:"type:text;default:'pending'" json:"status"`
	ExpiresAt      time.Time   `gorm:"type:timestamptz;not null" json:"expires_at"`
	PaidAt         *time.Time  `gorm:"type:timestamptz" json:"paid_at"`
	PaymentMethod  string      `gorm:"type:text" json:"payment_method"`
	PaymentChannel string      `gorm:"type:text" json:"payment_channel"`
	CreatedAt      time.Time   `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`

	// Relations
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template Template `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}
