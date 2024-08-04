package converter

import (
	"customer-service-backend/internal/entity"
	"customer-service-backend/internal/models"
)

func UserToResponse(user *entity.User) *models.UserResponse {
	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
}
