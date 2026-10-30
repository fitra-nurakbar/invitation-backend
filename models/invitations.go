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

// ─── Detail Structs ───────────────────────────────────────────

type BrideGroom struct {
	Name         string `json:"name"`
	PlaceOfBirth string `json:"place_of_birth"`
	Parent       string `json:"parent"`
	Photo        string `json:"photo"`
}

type Quotes struct {
	Quote       string `json:"quote"`
	Attribution string `json:"attribution"`
}

type WeddingEvent struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Place   string `json:"place"`
	Address string `json:"address"`
	MapsURL string `json:"maps_url"`
}

type WeddingGift struct {
	Platform  string `json:"platform"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	QRISImage string `json:"qris_image,omitempty"`
}

type InvitationDetailInput struct {
	CoverImage        string                 `json:"cover_image,omitempty"`
	Bride             BrideGroom             `json:"bride"`
	Groom             BrideGroom             `json:"groom"`
	Quotes            Quotes                 `json:"quotes"`
	WeddingEvent      []WeddingEvent         `json:"wedding_event"`
	StreamingPlatform string                 `json:"streaming_platform"`
	Gallery           map[string]interface{} `json:"gallery"`
	LoveStory         map[string]interface{} `json:"love_story"`
	WeddingGift       []WeddingGift          `json:"wedding_gift"`
	Music             string                 `json:"music"`
}

// ─── Model ────────────────────────────────────────────────────

type Invitation struct {
	ID              uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	TemplateID      uuid.UUID        `gorm:"type:uuid;not null" json:"template_id"`
	Slug            string           `gorm:"type:text;uniqueIndex;not null" json:"slug"`
	EventDate       time.Time        `gorm:"type:date" json:"event_date"`
	Status          InvitationStatus `gorm:"type:text;default:'draft'" json:"status"`
	ExpiresAt       time.Time        `gorm:"type:timestamptz" json:"expires_at"`
	Detail          datatypes.JSON   `gorm:"type:jsonb" json:"detail"`
	Layout          datatypes.JSON   `gorm:"type:jsonb;not null;default:'{}'::jsonb" json:"layout"`
	DraftLayout     datatypes.JSON   `gorm:"column:draft_layout;type:jsonb;not null;default:'{}'::jsonb" json:"draft_layout,omitempty"`
	PublishedLayout datatypes.JSON   `gorm:"column:published_layout;type:jsonb;not null;default:'{}'::jsonb" json:"published_layout,omitempty"`
	PublishedAt     *time.Time       `gorm:"column:published_at" json:"published_at,omitempty"`
	LayoutVer       int              `gorm:"column:layout_version;type:integer;not null;default:1" json:"layout_version"`
	CreatedAt       time.Time        `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relations
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template Template  `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Messages []Message `gorm:"foreignKey:InvitationID" json:"messages,omitempty"`
}
