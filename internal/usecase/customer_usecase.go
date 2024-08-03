package usecase

import (
	"customer-service-backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomerUsecase struct {
	DB                 *gorm.DB
	Logger             *zap.Logger
	CustomerRepository *repository.CustomerRepository
}

func NewCustomerUserCase(db *gorm.DB, logger *zap.Logger, customerRepository *repository.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{
		DB:                 db,
		Logger:             logger,
		CustomerRepository: customerRepository,
	}
}
