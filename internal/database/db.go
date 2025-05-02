package database

import (
	"Medods/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(fmt.Sprintf("Failed to auto-migrate database: %s", err))
	}
}
