package server

import (
	"log"
	"os"
	"user-auth/user/handlers"
	"user-auth/user/repositories"
	"user-auth/user/usecases"
	"user-auth/utils"

	"user-auth/db"

	messageWorker "github.com/moronimotta/message-worker-module"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQServer struct {
	db            db.Database
	connectionUrl string
	queueName     string
	topicName     string
}

func NewRabbitMQServer(db db.Database) *RabbitMQServer {
	utils.InitLogging()

	return &RabbitMQServer{
		db:            db,
		connectionUrl: os.Getenv("RABBITMQ_URL"),
		queueName:     os.Getenv("RABBITMQ_QUEUE_NAME"),
		topicName:     os.Getenv("RABBITMQ_TOPIC_NAME"),
	}
}
func (s *RabbitMQServer) Start() {
	// Setup repositories and handler
	userPostgresRepository := repositories.NewUserPostgresRepository(s.db)
	rabbitMqUserRepo := usecases.NewUserEventUsecase(userPostgresRepository)
	rabbitMqHandler := handlers.NewUserRabbitMQHandler(*rabbitMqUserRepo)

	// CONSUMER
	var consumerInput messageWorker.Consumer
	consumerInput.ConnectionURL = os.Getenv("RABBITMQ_URL")
	consumerInput.QueueName = s.queueName
	consumerInput.TopicName = s.topicName

	messageWorker.Listen(consumerInput,
		func(msg amqp.Delivery) {
			err := rabbitMqHandler.Repo.EventBus(string(msg.Body))
			if err != nil {
				log.Printf("EventBus error: %v", err)
			}
		},
	)

}
