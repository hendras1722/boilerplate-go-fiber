package database

import (
	"fmt"
	"log"

	"github.com/username/msa-boilerplate-go/config"
	"github.com/username/msa-boilerplate-go/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migration
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Printf("Error auto migrating database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
