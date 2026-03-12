package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"index;not null"`
	User      *User     `gorm:"constraint:OnDelete:CASCADE;"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false;not null"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
