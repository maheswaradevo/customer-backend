package repository

import (
	"customer-service-backend/internal/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TenorRepository struct {
	Repository[entity.Tenor]
	logger *zap.Logger
}

func NewTenorRepository(logger *zap.Logger) *TenorRepository {
	return &TenorRepository{logger: logger}
}

func (r *TenorRepository) Get(db *gorm.DB, id uint64) (*entity.Tenor, error) {
	var tenor *entity.Tenor

	if err := db.Where("id = ?", id).Find(&tenor).Error; err != nil {
		r.logger.Error("failed to get credit limit: ", zap.Error(err))
		return nil, err
	}

	return tenor, nil
}

func (r *TenorRepository) GetAll(db *gorm.DB) (*[]entity.Tenor, int64, error) {
	var result []entity.Tenor
	var count int64

	query := db.Model(&entity.Tenor{})

	err := query.Count(&count).Error

	if err != nil {
		r.logger.Error("failed to get all tenor: ", zap.Error(err))
		return nil, 0, err
	}

	err = query.Find(&result).Error

	if err != nil {
		r.logger.Error("failed to get all tenor: ", zap.Error(err))
		return nil, 0, err
	}

	return &result, count, nil
}
