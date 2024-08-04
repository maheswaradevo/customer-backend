package usecase

import (
	"customer-service-backend/internal/common/constants"
	"customer-service-backend/internal/common/helpers"
	"customer-service-backend/internal/models"
	"customer-service-backend/internal/models/converter"
	"customer-service-backend/internal/repository"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB                 *gorm.DB
	Log                *zap.Logger
	AuthRepository     *repository.AuthRepository
	CustomerRepository *repository.CustomerRepository
}

func NewUserUseCase(db *gorm.DB, logger *zap.Logger, authRepository *repository.AuthRepository, customerRepository *repository.CustomerRepository) *AuthUseCase {
	return &AuthUseCase{
		DB:                 db,
		Log:                logger,
		AuthRepository:     authRepository,
		CustomerRepository: customerRepository,
	}
}

func (u *AuthUseCase) Register(ctx echo.Context, request *models.CustomerRegisterRequest) (*models.UserResponse, error) {
	tx := u.DB.WithContext(helpers.Context(ctx)).Begin()
	defer tx.Rollback()

	userGet, err := u.AuthRepository.Get(tx, models.UserGetRequest{
		Username: request.Email,
	})

	if err != nil && err != gorm.ErrRecordNotFound {
		u.Log.Error("failed to get user data: ", zap.Error(err))
		return nil, err
	}

	if userGet != nil {
		u.Log.Error("cannot use the same email")
		err = errors.New("email can't be same")
		return nil, err
	}

	if len(request.IdNumber) != constants.NIKLen {
		err = errors.New("NIK should be 16 digits")
		return nil, err
	}

	if _, err := strconv.Atoi(request.IdNumber); err != nil {
		return nil, errors.New("NIK should only contain digits")
	}

	err = helpers.ValidateIDNumber(request.IdNumber)
	if err != nil {
		return nil, err
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Error("failed to generate encrypted password: ", zap.Error(err))
		return nil, err
	}

	idEncrypted := helpers.EncryptData(request.IdNumber, helpers.KeyEncrypDecryptData)

	request.IdNumber = idEncrypted

	user, err := u.AuthRepository.Create(tx, &models.UserCreateRequest{
		Username: request.Username,
		Password: string(hashPass),
	})
	if err != nil {
		tx.Rollback()
		u.Log.Error("failed to create user: ", zap.Error(err))
		return nil, err
	}

	if request.BirthdayDate != "" {
		birthDayTime, err := time.Parse(constants.DateFormat, request.BirthdayDate)
		if err != nil {
			tx.Rollback()
			u.Log.Error("failed to parse birthday date: ", zap.Error(err))
			return nil, err
		}
		request.BirthdayDateTime = birthDayTime
	}

	_, err = u.CustomerRepository.Create(tx, &models.CustomerCreateRequest{
		IdNumber:         idEncrypted,
		FullName:         request.FullName,
		LegalName:        request.LegalName,
		BirthdayLoc:      request.BirthdayLoc,
		BirthdayDate:     request.BirthdayDate,
		Salary:           request.Salary,
		BirthdayDateTime: request.BirthdayDateTime,
	})
	if err != nil {
		tx.Rollback()
		u.Log.Error("failed to create customer: ", zap.Error(err))
		return nil, err
	}

	tx.Commit()

	result := converter.UserToResponse(user)

	return result, nil
}

func (u *AuthUseCase) Login(ctx echo.Context, data models.LoginRequest) (*models.LoginResponse, error) {
	tx := u.DB.WithContext(helpers.Context(ctx)).Begin()
	defer tx.Rollback()
	user, err := u.AuthRepository.Get(tx, models.UserGetRequest{
		Username: data.Username,
		Email:    data.Email,
	})
	if err != nil {
		u.Log.Error("failed to get user: ", zap.Error(err))
		return nil, err
	}

	errMismatchPassowrd := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if errMismatchPassowrd != nil {
		u.Log.Error("wrong password")
		err = errors.New("wrong password")
		return nil, err
	}

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}
	accessToken, refreshToken, expiredToken, expRefreshToken, err := helpers.GenerateTokenPair(claims)
	if err != nil {
		u.Log.Error("failed to generate token pair: ", zap.Error(err))
		return nil, err
	}

	result := models.LoginResponse{
		UserID:              user.ID,
		Email:               user.Email,
		Username:            user.Username,
		AccessToken:         *accessToken,
		RefreshToken:        *refreshToken,
		ExpiredToken:        *expiredToken,
		ExpiredRefreshToken: *expRefreshToken,
	}

	return &result, nil
}
