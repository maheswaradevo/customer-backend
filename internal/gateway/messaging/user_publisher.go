package messaging

import (
	"customer-service-backend/internal/common"
	"customer-service-backend/internal/models"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type UserPublisher struct {
	event *models.Events
	log   *zap.Logger
}

func (u *UserPublisher) PushUserData(data *models.UserEvent) (bool, error) {
	pyld, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	request := models.MessageRabbit{
		Exchange: common.UserDataExchange,
		Queue:    common.UserDataSent,
		Payload:  pyld,
	}

	if _, err := u.Publish(request); err != nil {
		u.log.Error("failed to publish data: ", zap.Error(err))
		return false, err
	}

	return true, nil
}

func (u *UserPublisher) Publish(data models.MessageRabbit) (bool, error) {
	if err := u.event.Publisher.Publish(data.Queue, &message.Message{
		Payload: data.Payload,
	}); err != nil {
		u.log.Error("failed to broadcast message to "+data.Queue, zap.Error(err))
		return false, err
	}

	u.log.Info("successfully broadcast message to " + data.Queue)
	return true, nil
}

func NewUserPublisher(events *models.Events, log *zap.Logger) *UserPublisher {
	return &UserPublisher{
		log:   log,
		event: events,
	}
}
