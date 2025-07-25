package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"user-auth/db"
	"user-auth/user/entities"
	"user-auth/user/repositories"
	"user-auth/user/usecases"
	usescases "user-auth/user/usecases"

	messageWorker "github.com/moronimotta/message-worker-module"
	"github.com/redis/go-redis/v9"
)

type RabbitMqHandler struct {
	usescases.UserUsecase
	redisClient *redis.Client
}

func NewRabbitMqHandler(db db.Database, redisClient *redis.Client) *RabbitMqHandler {
	repoInput := repositories.NewUserPostgresRepository(db)
	usecasesInput := usecases.NewUserUsecase(repoInput)
	return &RabbitMqHandler{
		UserUsecase: *usecasesInput,
		redisClient: redisClient,
	}
}

func (h *RabbitMqHandler) EventBus(event messageWorker.Event) error {

	userEntity := &entities.User{}

	dataMap, ok := event.Data.(map[string]interface{})
	if !ok {
		return errors.New("event data is not of type map[string]interface{}")
	}

	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return fmt.Errorf("failed to marshal data map: %w", err)
	}

	if err := json.Unmarshal(jsonData, userEntity); err != nil {
		return fmt.Errorf("failed to unmarshal into User: %w", err)
	}

	switch event.Event {
	case "user.updated":

		if err := h.UpdateUser(userEntity); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("unhandled event type: %s", event.Event)
	}
}

func (u *RabbitMqHandler) PublishMessage(topicName, eventName string, data map[string]string) error {

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
