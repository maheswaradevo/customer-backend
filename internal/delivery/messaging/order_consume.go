package messaging

import (
	"context"
	"customer-service-backend/internal/common"
	"customer-service-backend/internal/models"
	"customer-service-backend/internal/models/consumer"
	"customer-service-backend/internal/usecase"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreditLimitConsumer struct {
	UseCase *usecase.AuthUseCase
}

func NewCreditLimitConsumer(useCase *usecase.AuthUseCase) *CreditLimitConsumer {
	return &CreditLimitConsumer{UseCase: useCase}
}

func (c *CreditLimitConsumer) ConsumeCreditLimitRequest(db *gorm.DB, log *zap.Logger, events *models.Events) {
	msgs, err := events.SubGroup("request").Subscribe(context.TODO(), common.CreditLimitExchange)
	if err != nil {
		log.Error(err.Error())
	}

	go func(msgs <-chan *message.Message) {
		for msg := range msgs {
			data := &consumer.CreditLimitRequest{}
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				continue
			}

			_, err := c.UseCase.HandleCreditLimitRequest(data.CustomerID)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			msg.Ack()
			log.Info("success consume message topic: " + common.UserDataSent)
		}
	}(msgs)
}

func (c *CreditLimitConsumer) ConsumeUpdateCreditLimit(db *gorm.DB, log *zap.Logger, events *models.Events) {
	msgs, err := events.SubGroup("update").Subscribe(context.TODO(), common.CreditLimitExchange)
	if err != nil {
		log.Error(err.Error())
	}

	go func(msgs <-chan *message.Message) {
		for msg := range msgs {
			data := &consumer.CreditLimitUpdate{}
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				continue
			}

			_, err := c.UseCase.HandleUpdateFromOrder(data)
			if err != nil {
				log.Error("failed to update: ", zap.Error(err))
				continue
			}

			msg.Ack()
			log.Info("success consume message topic: " + common.UserDataSent)
		}
	}(msgs)
}
