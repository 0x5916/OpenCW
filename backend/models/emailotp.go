package models

import "gorm.io/gorm"

type EmailOTP struct {
	gorm.Model
	Email string `gorm:"not null"`
	Code  string `gorm:"not null"`
}

func (EmailOTP) TableName() string {
	return "email_otp"
}
