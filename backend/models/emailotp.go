package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailOTP struct {
	gorm.Model
	UserID    uint      `gorm:"index"`
	Email     string    `gorm:"not null;index"`
	Code      string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
	Verified  bool      `gorm:"default:false;not null"`
}

func (EmailOTP) TableName() string {
	return "email_otp"
}
