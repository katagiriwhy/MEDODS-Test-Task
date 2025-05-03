package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(256);not null;unique"`
	Password     string `gorm:"type:varchar(256);not null"`
	GUID         string `gorm:"type:varchar(256);unique"`
	RefreshToken string `gorm:"type:varchar(512);not null;default:''"`
}
