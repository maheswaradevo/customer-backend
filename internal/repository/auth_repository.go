package repository

import (
	"customer-service-backend/internal/entity"
	"customer-service-backend/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepository struct {
	Repository[entity.User]
	logger *zap.Logger
}

func NewAuthRepository(logger *zap.Logger) *AuthRepository {
	return &AuthRepository{logger: logger}
}

func (r *AuthRepository) Create(db *gorm.DB, user *models.UserCreateRequest) (*entity.User, error) {
	result := entity.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	err := db.Model(&entity.User{}).Create(&result).Error
	if err != nil {
		r.logger.Error("failed to create user data: ", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (r *AuthRepository) Get(db *gorm.DB, data models.UserGetRequest) (*entity.User, error) {
	var user entity.User
	tx := db.Model(&user)

	if data.Email != "" {
		tx = tx.Where("email = ?", data.Email)
	}
	if data.Username != "" {
		tx = tx.Where("username = ?", data.Username)
	}
	if data.ID != 0 {
		tx = tx.Where("id = ?", data.ID)
	}

	err := tx.First(&user).Error
	if err != nil {
		r.logger.Error("failed to get user: ", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) Update(db *gorm.DB, data models.UserUpdateRequest) error {
	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.User{}).Where("id = ?", data.ID).Updates(&models.UserUpdateRequest{
			Password: data.Password,
			Username: data.Username,
		}).Error; err != nil {
			r.logger.Error("failed to update user: ", zap.Error(err))
			return err
		}
		return nil
	})
	return nil
}

func (r *AuthRepository) Delete(db *gorm.DB, id uint64) error {
	if err := db.Delete(&entity.User{}, db.Where("id = ?", id)).Error; err != nil {
		r.logger.Error("failed to delete user: ", zap.Error(err))
		return err
	}
	return nil
}
