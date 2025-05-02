package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"size:255;not null;unique"`
	Password string `gorm:"size:255;not null"`
	GUID     string `gorm:"size:255;not null;unique"`
}
