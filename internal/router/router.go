package router

import (
	"github.com/Verano-20/stage-zero/internal/config"
	"github.com/Verano-20/stage-zero/internal/container"
	"github.com/Verano-20/stage-zero/internal/controller"
	"github.com/Verano-20/stage-zero/internal/logger"
	"github.com/Verano-20/stage-zero/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func InitRouter(container *container.Container) *gin.Engine {
	config := config.Get()
	log := logger.Get()

	log.Info("Configuring router...")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(config.ServiceName))
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.MetricsMiddleware())

	authMiddleware := middleware.NewAuthMiddleware(config.GetJwtSecret(), container.UserRepository)

	router.GET("/health", controller.GetHealth)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth
	authController := container.AuthController
	auth := router.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	// Simple
	simpleController := container.SimpleController
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
