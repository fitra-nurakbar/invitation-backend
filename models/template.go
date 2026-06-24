package models

import "github.com/google/uuid"

type Template struct {
    ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Name             string    `gorm:"type:text;not null" json:"name"`
    Price            int       `gorm:"type:integer" json:"price"`
    IsActive         bool      `gorm:"type:boolean;default:true" json:"is_active"`
    OrderDeadlineDays int      `gorm:"type:integer" json:"order_deadline_days"`
    ActiveDaysAfter  int       `gorm:"type:integer" json:"active_days_after"`
}