package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Template struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Slug              string         `gorm:"type:text;uniqueIndex;not null" json:"slug"`
	Name              string         `gorm:"type:text;not null" json:"name"`
	Description       string         `gorm:"type:text;not null;default:''" json:"description"`
	Thumbnail         string         `gorm:"type:text;not null;default:''" json:"thumbnail"`
	PreviewURL        string         `gorm:"type:text;not null;default:''" json:"preview_url"`
	Category          string         `gorm:"type:text;not null;default:'wedding'" json:"category"`
	Price             int            `gorm:"type:integer" json:"price"`
	IsPremium         bool           `gorm:"type:boolean;not null;default:false" json:"is_premium"`
	IsActive          bool           `gorm:"type:boolean;default:true" json:"is_active"`
	OrderDeadlineDays int            `gorm:"type:integer" json:"order_deadline_days"`
	ActiveDaysAfter   int            `gorm:"type:integer" json:"active_days_after"`
	ThemeConfig       datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'::jsonb" json:"theme_config"`
	BlockConfig       datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'::jsonb" json:"block_config"`
	Version           int            `gorm:"type:integer;not null;default:1" json:"version"`
}
