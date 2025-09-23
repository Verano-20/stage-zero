package container

import (
	"github.com/Verano-20/go-crud/internal/container"
	"github.com/Verano-20/go-crud/internal/model"
	mocks "github.com/Verano-20/go-crud/test/mocks/repository"
)

func NewContainerWithMockRepositories() *container.Container {
	userRepository := mocks.NewMockUserRepository(map[uint]*model.User{})
	// simpleRepository := mocks.NewMockSimpleRepository() TODO MOCK SIMPLE REPO
	return container.NewContainerWithInterfaces(userRepository, nil)
}
