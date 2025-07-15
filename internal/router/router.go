package router

import (
	"github.com/Verano-20/go-crud/internal/controller"
	"github.com/Verano-20/go-crud/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitializeRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())

	authMiddleware := middleware.NewAuthMiddleware(db)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", controller.GetHealth)

	// Auth
	authController := controller.NewAuthController(db)
	auth := router.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	// Simple
	simpleController := controller.NewSimpleController(db)
	simples := router.Group("/simple", authMiddleware.AuthenticateRequest)
	{
		simples.POST("/", simpleController.Create)
		simples.GET("/", simpleController.GetAll)
		simples.GET("/:id", simpleController.GetByID)
		simples.PUT("/:id", simpleController.Update)
		simples.DELETE("/:id", simpleController.Delete)
	}

	return router
}
