package container

import (
	"github.com/Verano-20/stage-zero/internal/controller"
	"github.com/Verano-20/stage-zero/internal/repository"
	"github.com/Verano-20/stage-zero/internal/service"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB

	// Repositories
	UserRepository   repository.UserRepository
	SimpleRepository repository.SimpleRepository

	// Services
	UserService   service.UserService
	AuthService   service.AuthService
	SimpleService service.SimpleService

	// Controllers
	AuthController   *controller.AuthController
	SimpleController *controller.SimpleController
}

func NewContainerWithDB(db *gorm.DB) *Container {
	userRepository := repository.NewUserRepository(db)
	simpleRepository := repository.NewSimpleRepository(db)

	container := NewContainerWithInterfaces(userRepository, simpleRepository)
	container.DB = db
	return container
}

func NewContainerWithInterfaces(userRepository repository.UserRepository, simpleRepository repository.SimpleRepository) *Container {
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userService)
	simpleService := service.NewSimpleService(simpleRepository)

	authController := controller.NewAuthController(userService, authService)
	simpleController := controller.NewSimpleController(simpleService)

	return &Container{
		UserRepository:   userRepository,
		SimpleRepository: simpleRepository,
		UserService:      userService,
		AuthService:      authService,
		SimpleService:    simpleService,
		AuthController:   authController,
		SimpleController: simpleController,
	}
}
