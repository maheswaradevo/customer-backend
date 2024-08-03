package repository

import (
	"customer-service-backend/internal/entity"
	"customer-service-backend/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Repository[entity.Customer]
	logger *zap.Logger
}

func NewCustomerRepository(logger *zap.Logger) *CustomerRepository {
	return &CustomerRepository{logger: logger}
}

func (r *CustomerRepository) Create(db *gorm.DB, customer *models.CustomerCreateRequest) (*entity.Customer, error) {
	result := entity.Customer{
		IdNumber:     customer.IdNumber,
		FullName:     customer.FullName,
		LegalName:    customer.LegalName,
		BirthdayLoc:  customer.BirthdayLoc,
		BirthdayDate: customer.BirthdayDateTime,
		Salary:       customer.Salary,
		UpdatedAt:    nil,
	}
	err := db.Create(&result).Error
	if err != nil {
		r.logger.Error("failed to create customer data: ", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (r *CustomerRepository) Get(db *gorm.DB, id uint64) (*entity.Customer, error) {
	var customer *entity.Customer

	if err := db.Where("id = ?", id).Find(&customer).Error; err != nil {
		r.logger.Error("failed to get customer: ", zap.Error(err))
		return nil, err
	}

	return customer, nil
}

func (r *CustomerRepository) Update(db *gorm.DB, data models.CustomerUpdateRequest) error {
	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Customer{}).Where("id = ?", data.ID).Updates(&models.CustomerUpdateRequest{
			IdNumber:     data.IdNumber,
			FullName:     data.FullName,
			LegalName:    data.LegalName,
			BirthdayDate: data.BirthdayDate,
			BirthdayLoc:  data.BirthdayLoc,
			IdPicUrl:     data.IdPicUrl,
			SelfPicUrl:   data.SelfPicUrl,
		}).Error; err != nil {
			r.logger.Error("failed to update customer: ", zap.Error(err))
			return err
		}
		return nil
	})
	return nil
}

func (r *CustomerRepository) Delete(db *gorm.DB, id uint64) error {
	if err := db.Delete(&entity.Customer{}, db.Where("id = ?", id)).Error; err != nil {
		r.logger.Error("failed to delete customer data: ", zap.Error(err))
		return err
	}
	return nil
}
