package container

import (
	"github.com/Verano-20/go-crud/internal/container"
	"github.com/Verano-20/go-crud/internal/model"
	mocks "github.com/Verano-20/go-crud/test/mocks/repository"
)

var globalContainer *container.Container

func GetTestContainer() *container.Container {
	if globalContainer == nil {
		InitContainerWithMockRepositories()
	}
	return globalContainer
}

func InitContainerWithMockRepositories() {
	userRepository := mocks.NewMockUserRepository(map[uint]*model.User{})
	simpleRepository := mocks.NewMockSimpleRepository(map[uint]*model.Simple{})
	globalContainer = container.NewContainerWithInterfaces(userRepository, simpleRepository)
}
