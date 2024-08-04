package repository

import (
	"customer-service-backend/internal/entity"
	"customer-service-backend/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreditLimitRepository struct {
	Repository[entity.CreditLimit]
	logger *zap.Logger
}

func NewCreditLimitRepository(logger *zap.Logger) *CreditLimitRepository {
	return &CreditLimitRepository{logger: logger}
}

func (r *CreditLimitRepository) Create(db *gorm.DB, data *models.CreditLimitCreateRequest) (*entity.CreditLimit, error) {
	result := entity.CreditLimit{
		CustomerID:  data.CustomerID,
		CreditLimit: data.CreditLimit,
		TenorID:     uint(data.TenorID),
		StartDate:   data.StartDateTime,
		EndDate:     &data.EndDateTime,
	}
	err := db.Create(&result).Error
	if err != nil {
		r.logger.Error("failed to create customer data: ", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (r *CreditLimitRepository) Get(db *gorm.DB, id uint64) (*entity.CreditLimit, error) {
	var creditLimit *entity.CreditLimit

	if err := db.Where("id = ?", id).Find(&creditLimit).Error; err != nil {
		r.logger.Error("failed to get credit limit: ", zap.Error(err))
		return nil, err
	}

	return creditLimit, nil
}

func (r *CreditLimitRepository) Update(db *gorm.DB, data models.CreditLimitUpdateRequest) error {
	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.CreditLimit{}).Where("id = ?", data.ID).Updates(&models.CreditLimitUpdateRequest{
			CreditLimit:   data.CreditLimit,
			StartDateTime: data.StartDateTime,
			EndDateTime:   data.EndDateTime,
		}).Error; err != nil {
			r.logger.Error("failed to update credit limit: ", zap.Error(err))
			return err
		}
		return nil
	})
	return nil
}

func (r *CreditLimitRepository) Delete(db *gorm.DB, id uint64) error {
	if err := db.Delete(&entity.CreditLimit{}, db.Where("id = ?", id)).Error; err != nil {
		r.logger.Error("failed to delete credit limit data: ", zap.Error(err))
		return err
	}
	return nil
}
