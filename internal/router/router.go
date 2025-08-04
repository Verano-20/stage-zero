package router

import (
	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/controller"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	config := config.Get()
	log := logger.Get()

	log.Info("Configuring router...")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(config.ServiceName))
	router.Use(middleware.LoggingMiddleware())

	authMiddleware := middleware.NewAuthMiddleware(config.GetJwtSecret(), db)

	router.GET("/health", controller.GetHealth)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	log.Info("Router configured")
	return router
}
