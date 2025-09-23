package service

import (
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SimpleService struct {
	SimpleRepository repository.SimpleRepository
}

func NewSimpleService(simpleRepository repository.SimpleRepository) *SimpleService {
	return &SimpleService{SimpleRepository: simpleRepository}
}

func (s *SimpleService) CreateSimple(ctx *gin.Context, simpleForm model.SimpleForm) (*model.Simple, error) {
	log := logger.GetFromContext(ctx)

	log.Debug("Creating Simple...", zap.Object("simple", &simpleForm))

	simple, err := s.SimpleRepository.Create(ctx, simpleForm.ToModel())
	if err != nil {
		log.Error("Failed to create Simple",
			zap.Object("simple", &simpleForm),
			zap.Error(err))
		return nil, err
	}

	log.Debug("Simple created successfully", zap.Object("simple", simple))
	return simple, nil
}

func (s *SimpleService) GetAllSimples(ctx *gin.Context) (model.Simples, error) {
	log := logger.GetFromContext(ctx)

	log.Debug("Retrieving all Simples")

	simples, err := s.SimpleRepository.GetAll(ctx)
	if err != nil {
		log.Error("Failed to retrieve Simples", zap.Error(err))
		return nil, err
	}

	log.Debug("Simples retrieved successfully", zap.Int("count", len(simples)))
	return simples, nil
}

func (s *SimpleService) GetSimpleByID(ctx *gin.Context, id uint64) (*model.Simple, error) {
	log := logger.GetFromContext(ctx)

	log.Debug("Retrieving Simple by ID", zap.Uint64("id", id))

	simple, err := s.SimpleRepository.GetByID(ctx, uint(id))
	if err != nil {
		log.Warn("Simple not found", zap.Uint64("id", id), zap.Error(err))
		return nil, err
	}

	log.Debug("Simple retrieved successfully", zap.Object("simple", simple))
	return simple, nil
}

func (s *SimpleService) UpdateSimple(ctx *gin.Context, existingSimple *model.Simple, simpleForm model.SimpleForm) (*model.Simple, error) {
	log := logger.GetFromContext(ctx)

	log.Debug("Updating Simple",
		zap.Object("existing", existingSimple),
		zap.Object("update", &simpleForm))

	existingSimple.Name = simpleForm.Name

	simple, err := s.SimpleRepository.Update(ctx, existingSimple)
	if err != nil {
		log.Error("Failed to update Simple",
			zap.Object("existing", existingSimple),
			zap.Object("update", &simpleForm),
			zap.Error(err))
		return nil, err
	}

	log.Debug("Simple updated successfully", zap.Object("updated", simple))
	return simple, nil
}

func (s *SimpleService) DeleteSimple(ctx *gin.Context, existingSimple *model.Simple) error {
	log := logger.GetFromContext(ctx)

	log.Debug("Deleting Simple", zap.Object("simple", existingSimple))

	err := s.SimpleRepository.Delete(ctx, existingSimple.ID)
	if err != nil {
		log.Error("Failed to delete Simple",
			zap.Object("simple", existingSimple),
			zap.Error(err))
		return err
	}

	log.Debug("Simple deleted successfully", zap.Object("simple", existingSimple))
	return nil
}
