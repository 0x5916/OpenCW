package models

import "gorm.io/gorm"

type Progress struct {
	gorm.Model
	UserID   uint    `gorm:"index;not null"`
	User     *User   `gorm:"constraint:OnDelete:CASCADE;"`
	Lesson   int     `gorm:"not null"`
	CharWPM  int     `gorm:"not null"`
	EffWPM   int     `gorm:"not null"`
	Accuracy float64 `gorm:"not null"`
}

func (Progress) TableName() string {
	return "lesson_progress"
}
