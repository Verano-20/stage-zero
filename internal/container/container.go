package container

import (
	"github.com/Verano-20/go-crud/internal/controller"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/Verano-20/go-crud/internal/service"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB

	// Repositories
	UserRepository   repository.UserRepository
	SimpleRepository repository.SimpleRepository

	// Services
	UserService   *service.UserService
	AuthService   *service.AuthService
	SimpleService *service.SimpleService

	// Controllers
	AuthController   *controller.AuthController
	SimpleController *controller.SimpleController
}

func NewContainer(db *gorm.DB) *Container {
	userRepo := repository.NewUserRepository(db)
	simpleRepo := repository.NewSimpleRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userService)
	simpleService := service.NewSimpleService(simpleRepo)

	authController := controller.NewAuthController(userService, authService)
	simpleController := controller.NewSimpleController(simpleService)

	return &Container{
		DB:               db,
		UserRepository:   userRepo,
		SimpleRepository: simpleRepo,
		UserService:      userService,
		AuthService:      authService,
		SimpleService:    simpleService,
		AuthController:   authController,
		SimpleController: simpleController,
	}
}
