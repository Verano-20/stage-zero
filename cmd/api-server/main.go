package main

import (
	"log"

	"github.com/Verano-20/go-crud/internal/initializer"
	"github.com/Verano-20/go-crud/internal/router"

	_ "github.com/Verano-20/go-crud/docs"
)

// @title           Go-CRUD API
// @version         1.0
// @description     This is a CRUD API for a simple application.

// @host      localhost:8080
// @BasePath  /
func main() {
	db := initializer.InitializeDatabase()
	router := router.InitializeRouter(db)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
