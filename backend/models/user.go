package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	CallSign     string        `gorm:"uniqueIndex;"`
	Username     string        `gorm:"uniqueIndex;not null"`
	Email        string        `gorm:"uniqueIndex;not null"`
	Password     string        `gorm:"not null"`
	CWSettings   *CWSettings   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
	PageSettings *PageSettings `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (User) TableName() string {
	return "users"
}
