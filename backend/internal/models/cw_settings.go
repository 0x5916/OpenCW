package models

import (
	"github.com/google/uuid"
)

type CWSettings struct {
	Base
	UserID     uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	User       *User     `gorm:"constraint:OnDelete:CASCADE;"`
	CharWPM    int       `gorm:"not null"`
	EffWPM     int       `gorm:"not null"`
	Freq       int       `gorm:"not null"`
	StartDelay float64   `gorm:"not null"`
}

func (CWSettings) TableName() string {
	return "cw_settings"
}

func GetDefaultCWSettings() CWSettings {
	return CWSettings{
		CharWPM:    20,
		EffWPM:     12,
		Freq:       600,
		StartDelay: 0.5,
	}
}
