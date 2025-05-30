package usecases

import (
	"os"

	messageWorker "github.com/moronimotta/message-worker-module"
)

func NewRabbitUsecase(repo UserEventUsecase) *RabbitUsecase {
	return &RabbitUsecase{
		repo,
	}
}

type RabbitUsecase struct {
	UserEventUsecase
}

func (u *RabbitUsecase) PublishMessage(topicName, eventName string, data map[string]string) error {

	input := messageWorker.Publisher{}
	input.ConnectionURL = os.Getenv("RABBITMQ_URL")
	input.TopicName = topicName

	messageInput := messageWorker.Event{
		Event: eventName,
		Data:  data,
	}

	messageWorker.SendMessage(input, messageInput)
	return nil
}
