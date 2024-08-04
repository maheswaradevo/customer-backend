package route

import (
	"customer-service-backend/internal/config"
	"customer-service-backend/internal/delivery/http"
	consumer "customer-service-backend/internal/delivery/messaging"
	"customer-service-backend/internal/gateway/messaging"
	"customer-service-backend/internal/models"
	"customer-service-backend/internal/repository"
	"customer-service-backend/internal/usecase"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteConfig struct {
	App            *echo.Echo
	AuthController *http.AuthController
}

type BootstrapConfig struct {
	DB           *gorm.DB
	Redis        *redis.Client
	App          *echo.Echo
	Log          *zap.Logger
	Config       *config.Config
	Events       models.Events
	RabbitMQConn *amqp.Connection
	RabbitMQChan *amqp.Channel
	RabbitMQQuit chan bool
}

func Bootstrap(config *BootstrapConfig) {
	// Setup Repositories
	authRepository := repository.NewAuthRepository(config.Log)
	customerRepository := repository.NewCustomerRepository(config.Log)
	tenorRepository := repository.NewTenorRepository(config.Log)
	creditLimitRepository := repository.NewCreditLimitRepository(config.DB, config.Log)

	// Setup PubSub
	userMessaging := messaging.NewUserPublisher(&config.Events, config.Log)
	creditLimitMessaging := messaging.NewCreditLimitPublisher(&config.Events, config.Log)

	// Setup usecases
	authUseCase := usecase.NewUserUseCase(config.DB, config.Log, authRepository, customerRepository, tenorRepository, creditLimitRepository, userMessaging, creditLimitMessaging, config.Redis)

	creditLimitConsumer := consumer.NewCreditLimitConsumer(authUseCase)
	// Setup Controller
	authController := http.NewAuthController(authUseCase, config.Log)

	routeConfig := RouteConfig{
		App:            config.App,
		AuthController: authController,
	}

	time.AfterFunc(2*time.Second, func() {
		creditLimitConsumer.ConsumeCreditLimitRequest(config.DB, config.Log, &config.Events)
		creditLimitConsumer.ConsumeUpdateCreditLimit(config.DB, config.Log, &config.Events)
	})

	routeConfig.Setup()
}

func (c *RouteConfig) Setup() {
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	// Setup endpoint
	c.App.POST("/api/users", c.AuthController.Register)
	c.App.POST("/api/user/login", c.AuthController.Login)
}
