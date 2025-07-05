package repository

import (
	"github.com/Verano-20/go-crud/internal/model"
	"gorm.io/gorm"
)

type SimpleRepository struct {
	DB *gorm.DB
}

func NewSimpleRepository(db *gorm.DB) *SimpleRepository {
	return &SimpleRepository{DB: db}
}

func (r *SimpleRepository) Create(simple *model.Simple) (*model.Simple, error) {
	if err := r.DB.Create(&simple).Error; err != nil {
		return nil, err
	}

	return simple, nil
}

func (r *SimpleRepository) GetAll() (model.Simples, error) {
	var simples model.Simples
	if err := r.DB.Find(&simples).Error; err != nil {
		return nil, err
	}

	return simples, nil
}

func (r *SimpleRepository) GetByID(id uint) (*model.Simple, error) {
	simple := &model.Simple{}
	if err := r.DB.First(&simple, id).Error; err != nil {
		return nil, err
	}

	return simple, nil
}

func (r *SimpleRepository) Update(simple *model.Simple) (*model.Simple, error) {
	if err := r.DB.Save(&simple).Error; err != nil {
		return nil, err
	}

	return simple, nil
}

func (r *SimpleRepository) Delete(id uint) error {
	if err := r.DB.Delete(&model.Simple{}, id).Error; err != nil {
		return err
	}

	return nil
}
