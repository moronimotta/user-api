package server

import (
	"log"
	"user-auth/user/handlers"
	"user-auth/user/repositories"
	"user-auth/user/usecases"
	"user-auth/utils"

	"user-auth/db"

	messageWorker "github.com/moronimotta/message-worker-module"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQServer struct {
	db           db.Database
	rabbitmqUrl  string
	queueName    string
	exchangeName string
}

func NewRabbitMQServer(db db.Database, url, queueName, exchangeName string) *RabbitMQServer {
	utils.InitLogging()

	return &RabbitMQServer{
		db:           db,
		rabbitmqUrl:  url,
		queueName:    queueName,
		exchangeName: exchangeName,
	}
}
func (s *RabbitMQServer) Start() {
	conn, err := amqp.Dial(s.rabbitmqUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Setup repositories and handler
	userPostgresRepository := repositories.NewUserPostgresRepository(s.db)
	rabbitMqUserRepo := usecases.NewUserEventUsecase(userPostgresRepository)
	rabbitMqHandler := handlers.NewUserRabbitMQHandler(*rabbitMqUserRepo)

	// CONSUMER
	messageWorker.Worker(conn, s.queueName, s.exchangeName,
		func(msg amqp.Delivery) {
			err := rabbitMqHandler.Repo.EventBus(string(msg.Body))
			if err != nil {
				log.Printf("EventBus error: %v", err)
			}
		},
	)

}
