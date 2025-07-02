package main

import (
	"log"

	"github.com/Verano-20/go-crud/internal/initializers"
	"github.com/Verano-20/go-crud/internal/router"
)

func main() {
	db := initializers.InitializeDatabase()
	router := router.InitializeRouter(db)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
