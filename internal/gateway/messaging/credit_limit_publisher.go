package messaging

import (
	"customer-service-backend/internal/common"
	"customer-service-backend/internal/models"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type CreditLimitPublisher struct {
	event *models.Events
	log   *zap.Logger
}

func (c *CreditLimitPublisher) PushCreditLimitData(data *[]models.CreditLimitEvent) (bool, error) {
	pyld, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	request := models.MessageRabbit{
		Exchange: common.CreditLimitExchange,
		Queue:    common.CreditLimitDataQueue,
		Payload:  pyld,
	}

	if _, err := c.Publish(request); err != nil {
		c.log.Error("failed to publish data: ", zap.Error(err))
		return false, err
	}

	return true, nil
}

func (c *CreditLimitPublisher) Publish(data models.MessageRabbit) (bool, error) {
	if err := c.event.Publisher.Publish(data.Queue, &message.Message{
		Payload: data.Payload,
	}); err != nil {
		c.log.Error("failed to broadcast message to "+data.Queue, zap.Error(err))
		return false, err
	}

	c.log.Info("successfully broadcast message to " + data.Queue)
	return true, nil
}

func NewCreditLimitPublisher(events *models.Events, log *zap.Logger) *CreditLimitPublisher {
	return &CreditLimitPublisher{
		log:   log,
		event: events,
	}
}
