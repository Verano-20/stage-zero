package main

import (
	"log"

	"github.com/Verano-20/go-crud/internal/handlers"
	"github.com/Verano-20/go-crud/internal/initializers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db := initializers.InitializeDatabase()

	// Initialize router
	router := gin.Default()

	// Initialize handlers
	simpleHandler := handlers.CreateSimpleHandler(db)

	// Routes
	simples := router.Group("/simple")
	{
		simples.POST("/", simpleHandler.Create)
		simples.GET("/", simpleHandler.GetAll)
		simples.GET("/:id", simpleHandler.GetByID)
		simples.PUT("/:id", simpleHandler.Update)
		simples.DELETE("/:id", simpleHandler.Delete)
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
