package initializers

import (
	"log"

	"github.com/Verano-20/go-crud/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	// Get database connection string
	dsn := config.GetDBConnectionString()

	// Initialize database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}
