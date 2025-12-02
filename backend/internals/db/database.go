package db

import (
	"aetrix/observer/internals/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDB() *gorm.DB {

	dsn := "host=localhost user=user password=password dbname=dbname port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	// Migrate Models
	err = db.AutoMigrate(&models.User{}, &models.Machine{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate models: %v", err)
	}
	// Seed initial data
	Seed(db)
	
	log.Println("✅ Connected to PostgreSQL & Migrated")
	return db
}
