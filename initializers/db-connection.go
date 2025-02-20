package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB instance
var DB *gorm.DB

func DBConnection() {
	// Get the connection string from environment variables
	connectionString := os.Getenv("CONNECTION_STRING")

	// Connect to the PostgreSQL database using GORM and the Postgres driver
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// If the connection is successful
	log.Println("Database connection established successfully")
}
