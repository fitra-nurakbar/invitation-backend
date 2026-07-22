package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"type:text" json:"name"`
	Password  *string   `gorm:"type:text" json:"-"`
	GoogleID  *string   `gorm:"type:text;uniqueIndex" json:"-"`
	Avatar    *string   `gorm:"type:text" json:"avatar,omitempty"`
	Role      UserRole  `gorm:"type:text;default:'user'" json:"role"`
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) HasPassword() bool {
	return u.Password != nil && *u.Password != ""
}

func (u *User) HasGoogleID() bool {
	return u.GoogleID != nil && *u.GoogleID != ""
}
