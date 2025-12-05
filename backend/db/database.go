package db

import (
	"log"
	"os"

	"election-pulse-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not set, ensure it is configured in .env")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	log.Println("Database connected successfully")

	// Auto Migrate
	err = DB.AutoMigrate(&models.Party{}, &models.SentimentSnapshot{})
	if err != nil {
		log.Printf("Failed to auto migrate: %v", err)
	}
}
