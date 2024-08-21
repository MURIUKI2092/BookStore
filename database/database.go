package database

import (
	"BookStore/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encapsulate the database connection and provides methods for interacting with the database
var DB *gorm.DB

func Connect() {
	var err error

	//Load environmenatal variables
	loaded := godotenv.Load()

	if loaded != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("databaseUser")
	dbPassword := os.Getenv("databasePassword")
	dbName := os.Getenv("databaseName")
	dbPort := os.Getenv("databasePort")
	dbHost := os.Getenv("databaseHost")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// handle migrations
	// AutoMigrate runs the migrations
	DB.AutoMigrate(&models.User{}, &models.Store{})

}
