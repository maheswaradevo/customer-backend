package usecase

import (
	"context"
	"customer-service-backend/internal/common/constants"
	"customer-service-backend/internal/common/helpers"
	"customer-service-backend/internal/gateway/messaging"
	"customer-service-backend/internal/models"
	"customer-service-backend/internal/models/converter"
	"customer-service-backend/internal/repository"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB                    *gorm.DB
	Log                   *zap.Logger
	AuthRepository        *repository.AuthRepository
	CustomerRepository    *repository.CustomerRepository
	TenorRepository       *repository.TenorRepository
	CreditLimitRepository *repository.CreditLimitRepository
	RedisClient           *redis.Client
	UserMessaging         *messaging.UserPublisher
	CreditLimitMessaging  *messaging.CreditLimitPublisher
}

func NewUserUseCase(db *gorm.DB, logger *zap.Logger, authRepository *repository.AuthRepository, customerRepository *repository.CustomerRepository, tenorRepository *repository.TenorRepository, creditLimitRepository *repository.CreditLimitRepository, userMessaging *messaging.UserPublisher, creditLimitMessaging *messaging.CreditLimitPublisher, redisClient *redis.Client) *AuthUseCase {
	return &AuthUseCase{
		DB:                    db,
		Log:                   logger,
		AuthRepository:        authRepository,
		CustomerRepository:    customerRepository,
		UserMessaging:         userMessaging,
		TenorRepository:       tenorRepository,
		CreditLimitRepository: creditLimitRepository,
		RedisClient:           redisClient,
		CreditLimitMessaging:  creditLimitMessaging,
	}
}

func (u *AuthUseCase) Register(ctx echo.Context, request *models.CustomerRegisterRequest) (*models.UserResponse, error) {
	tx := u.DB.WithContext(helpers.Context(ctx)).Begin()
	defer tx.Rollback()

	randGen := newRandomGenerator()

	userGet, err := u.AuthRepository.Get(tx, models.UserGetRequest{
		Email: request.Email,
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
		Email:    request.Email,
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

	tenors, _, err := u.TenorRepository.GetAll(tx)
	if err != nil {
		tx.Rollback()
		u.Log.Error("failed to get all tenor: ", zap.Error(err))
		return nil, err
	}

	var wg sync.WaitGroup
	var creditLimitError error

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range *tenors {
			var month int
			if v.TenorDescription == constants.OneMonth {
				month = 1
			} else if v.TenorDescription == constants.TwoMonth {
				month = 2
			} else if v.TenorDescription == constants.ThreeMonth {
				month = 3
			} else if v.TenorDescription == constants.SixMonth {
				month = 6
			}
			creditLimit := generateRandomCreditLimit(randGen, 1000000, 3000000)
			_, err := u.CreditLimitRepository.Create(tx, &models.CreditLimitCreateRequest{
				CustomerID:    user.ID,
				CreditLimit:   float64(creditLimit),
				TenorID:       v.ID,
				StartDateTime: time.Now(),
				EndDateTime:   time.Now().AddDate(0, month, 0),
			})
			if err != nil {
				creditLimitError = err
			}
		}
	}()

	wg.Wait()

	if creditLimitError != nil {
		tx.Rollback()
		u.Log.Error("failed to create credit limit: ", zap.Error(creditLimitError))
		return nil, creditLimitError
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

	userIdStr := strconv.FormatUint(user.ID, 10)

	ctxRedis := context.Background()

	err = u.RedisClient.Set(ctxRedis, fmt.Sprintf("%s:%s", constants.AccessTokenKey, userIdStr), *accessToken, time.Until(time.Now().Add(*expiredToken))).Err()
	if err != nil {
		u.Log.Error("failed to save access token to Redis: ", zap.Error(err))
		return nil, err
	}

	err = u.RedisClient.Set(ctxRedis, fmt.Sprintf("%s:%s", constants.RefreshTokenKey, userIdStr), *refreshToken, time.Until(time.Now().Add(*expRefreshToken))).Err()
	if err != nil {
		u.Log.Error("failed to save refresh token to Redis: ", zap.Error(err))
		return nil, err
	}

	go u.PublishUserLoginData(&result)

	return &result, nil
}

func (u *AuthUseCase) HandleCreditLimitRequest(customerID uint64) (bool, error) {
	creditLimit, _, err := u.CreditLimitRepository.GetAll(customerID)
	if err != nil {
		u.Log.Error(err.Error())
		return false, err
	}
	var eventData []models.CreditLimitEvent

	for _, v := range *creditLimit {
		eventData = append(eventData, models.CreditLimitEvent{
			ID:          v.ID,
			CustomerID:  v.CustomerID,
			CreditLimit: v.CreditLimit,
			StartDate:   v.StartDate,
			EndDate:     *v.EndDate,
		})
	}

	go u.PublishCreditLimitData(&eventData)
	return true, nil
}

func (u *AuthUseCase) PublishUserLoginData(data *models.LoginResponse) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			u.Log.Error("recovered from panic ", zap.Any("error", err))
		}
	}()

	return u.UserMessaging.PushUserData(&models.UserEvent{
		ID:                  data.UserID,
		Email:               data.Email,
		Username:            data.Username,
		AccessToken:         data.AccessToken,
		RefreshToken:        data.RefreshToken,
		ExpiredToken:        data.ExpiredToken,
		ExpiredRefreshToken: data.ExpiredRefreshToken,
	})
}

func (u *AuthUseCase) PublishCreditLimitData(data *[]models.CreditLimitEvent) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			u.Log.Error("recovered from panic ", zap.Any("error", err))
		}
	}()

	return u.CreditLimitMessaging.PushCreditLimitData(data)
}

func newRandomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func generateRandomCreditLimit(r *rand.Rand, min, max int) int {
	return r.Intn(max-min+1) + min
}
