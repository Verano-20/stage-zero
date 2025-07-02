package router

import (
	"github.com/Verano-20/go-crud/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitializeRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", handlers.GetHealth)

	// Simple
	simpleHandler := handlers.NewSimpleHandler(db)
	simples := router.Group("/simple")
	{
		simples.POST("/", simpleHandler.Create)
		simples.GET("/", simpleHandler.GetAll)
		simples.GET("/:id", simpleHandler.GetByID)
		simples.PUT("/:id", simpleHandler.Update)
		simples.DELETE("/:id", simpleHandler.Delete)
	}

	return router
}
