package router

import (
	controller "github.com/Verano-20/go-crud/internal/controller"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitializeRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", controller.GetHealth)

	// Auth
	userController := controller.NewUserController(db)
	users := router.Group("/auth")
	{
		users.POST("/signup", userController.SignUp)
		users.POST("/login", userController.Login)
	}

	// Simple
	simpleController := controller.NewSimpleController(db)
	simples := router.Group("/simple")
	{
		simples.POST("/", simpleController.Create)
		simples.GET("/", simpleController.GetAll)
		simples.GET("/:id", simpleController.GetByID)
		simples.PUT("/:id", simpleController.Update)
		simples.DELETE("/:id", simpleController.Delete)
	}

	return router
}
