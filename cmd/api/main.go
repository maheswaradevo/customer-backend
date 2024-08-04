package main

import (
	"customer-service-backend/internal/common/constants"
	"customer-service-backend/internal/config"
	"customer-service-backend/internal/delivery/http/route"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	config.Init()

	cfg := config.GetConfig()

	logger := config.NewLogger(cfg)
	db := config.NewDatabase(cfg, logger)
	app := config.NewEcho(cfg)
	rabbitMqClient, err := config.NewRabbitMQ(cfg.RabbitMqConfig)
	if err != nil {
		logger.Fatal("failed to start rabbitmq client: ", zap.Error(err))
	}
	rdb := config.NewRedis(cfg, logger)

	route.Bootstrap(&route.BootstrapConfig{
		DB:           db,
		Log:          logger,
		App:          app,
		Config:       cfg,
		RabbitMQConn: rabbitMqClient.Conn,
		RabbitMQChan: rabbitMqClient.Chann,
		RabbitMQQuit: rabbitMqClient.QuitChann,
		Events:       config.NewEvent(cfg),
		Redis:        rdb,
	})

	var address string
	if cfg.AppEnvironment == constants.LocalEnv {
		address = fmt.Sprintf("%s:%s", "localhost", cfg.Port)
	} else if cfg.AppEnvironment == constants.ProductionEnv {
		address = fmt.Sprintf("%s:%s", cfg.ProductionEnvironment, cfg.Port)
	}

	if err := app.Start(address); err != nil {
		logger.Fatal("failed to start server: ", zap.Error(err))
	}
}
