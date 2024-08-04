package repository

import (
	"customer-service-backend/internal/entity"
	"customer-service-backend/internal/models"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreditLimitRepository struct {
	Repository[entity.CreditLimit]
	logger *zap.Logger
	DB     *gorm.DB
}

func NewCreditLimitRepository(DB *gorm.DB, logger *zap.Logger) *CreditLimitRepository {
	return &CreditLimitRepository{
		DB:     DB,
		logger: logger,
	}
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

func (r *CreditLimitRepository) GetAll(customerID uint64) (*[]entity.CreditLimit, int64, error) {
	var result []entity.CreditLimit
	var count int64

	if r == nil {
		r.logger.Error("CreditLimitRepository is nil")
	}
	if r.logger == nil {
		// log a message using a fallback logger or panic
		fmt.Println("Logger is nil")
	}
	if r.DB == nil {
		// log a message using a fallback logger or panic
		fmt.Println("Database connection is nil")
	}

	query := r.DB.Model(&entity.CreditLimit{}).Where("customer_id = ?", customerID)

	err := query.Count(&count).Error
	if err != nil {
		r.logger.Error("failed to get all credit limit: ", zap.Error(err))
		return nil, 0, err
	}

	err = query.Find(&result).Error
	if err != nil {
		r.logger.Error("failed to get all credit limit: ", zap.Error(err))
		return nil, 0, err
	}

	return &result, count, nil
}
