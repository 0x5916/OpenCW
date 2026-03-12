package models

import "gorm.io/gorm"

type PageSettings struct {
	gorm.Model
	UserID uint   `gorm:"uniqueIndex;not null"`
	User   *User  `gorm:"constraint:OnDelete:CASCADE;"`
	Theme  string `gorm:"not null"`
	Lang   string `gorm:"not null"`
}

func (PageSettings) TableName() string {
	return "page_settings"
}

func GetDefaultPageSettings() PageSettings {
	return PageSettings{
		Theme: "auto",
		Lang:  "auto",
	}
}
