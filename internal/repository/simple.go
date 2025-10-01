package repository

import (
	"time"

	"github.com/Verano-20/stage-zero/internal/model"
	"github.com/Verano-20/stage-zero/internal/telemetry"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SimpleRepository interface {
	Create(ctx *gin.Context, simple *model.Simple) (*model.Simple, error)
	GetAll(ctx *gin.Context) (model.Simples, error)
	GetByID(ctx *gin.Context, id uint) (*model.Simple, error)
	Update(ctx *gin.Context, simple *model.Simple) (*model.Simple, error)
	Delete(ctx *gin.Context, id uint) error
}

type simpleRepository struct {
	DB *gorm.DB
}

var _ SimpleRepository = &simpleRepository{}

func NewSimpleRepository(db *gorm.DB) SimpleRepository {
	return &simpleRepository{DB: db}
}

func (r simpleRepository) Create(ctx *gin.Context, simple *model.Simple) (*model.Simple, error) {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	if err := r.DB.Create(&simple).Error; err != nil {
		return nil, err
	}

	metrics.RecordDBQuery(ctx, "create_simple", time.Since(start).Seconds())
	metrics.UpdateSimpleCount(ctx, 1)
	return simple, nil
}

func (r simpleRepository) GetAll(ctx *gin.Context) (model.Simples, error) {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	var simples model.Simples
	if err := r.DB.Find(&simples).Error; err != nil {
		return nil, err
	}

	metrics.RecordDBQuery(ctx, "get_all_simples", time.Since(start).Seconds())
	return simples, nil
}

func (r simpleRepository) GetByID(ctx *gin.Context, id uint) (*model.Simple, error) {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	simple := &model.Simple{}
	if err := r.DB.First(&simple, id).Error; err != nil {
		return nil, err
	}

	metrics.RecordDBQuery(ctx, "get_simple_by_id", time.Since(start).Seconds())
	return simple, nil
}

func (r simpleRepository) Update(ctx *gin.Context, simple *model.Simple) (*model.Simple, error) {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	if err := r.DB.Save(&simple).Error; err != nil {
		return nil, err
	}

	metrics.RecordDBQuery(ctx, "update_simple", time.Since(start).Seconds())
	return simple, nil
}

func (r simpleRepository) Delete(ctx *gin.Context, id uint) error {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	if err := r.DB.Delete(&model.Simple{}, id).Error; err != nil {
		return err
	}

	metrics.RecordDBQuery(ctx, "delete_simple", time.Since(start).Seconds())
	metrics.UpdateSimpleCount(ctx, -1)
	return nil
}
