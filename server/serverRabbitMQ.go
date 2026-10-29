package server

import (
	"os"
	"user-auth/user/handlers"
	"user-auth/utils"

	"user-auth/db"

	messageWorker "github.com/moronimotta/message-worker-module"
	"github.com/redis/go-redis/v9"
)

type RabbitMQServer struct {
	db            db.Database
	connectionUrl string
	queueName     string
	topicName     string
	redisClient   *redis.Client
}

func NewRabbitMQServer(db db.Database, redisClient *redis.Client) *RabbitMQServer {
	utils.InitLogging()

	return &RabbitMQServer{
		db:            db,
		connectionUrl: os.Getenv("RABBITMQ_URL"),
		queueName:     os.Getenv("RABBITMQ_QUEUE_NAME"),
		topicName:     os.Getenv("RABBITMQ_TOPIC_NAME"),
		redisClient:   redisClient,
	}
}
func (s *RabbitMQServer) Start() {
	// Setup repositories and handler
	rabbitMqHandler := handlers.NewRabbitMqHandler(s.db, s.redisClient)

	// CONSUMER
	var consumerInput messageWorker.Consumer
	consumerInput.ConnectionURL = os.Getenv("RABBITMQ_URL")
	consumerInput.QueueName = s.queueName
	consumerInput.TopicName = s.topicName

	messageWorker.Listen(consumerInput,
		func(event messageWorker.Event) error {
			return rabbitMqHandler.EventBus(event)
		},
	)

}
