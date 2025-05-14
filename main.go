package main

import (
	"log"
	"time"
	"user-auth/confs"
	"user-auth/db"
	"user-auth/server"
)

func main() {
	// load config
	config := confs.Config{}
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// connect to database Postgres
	database, err := db.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// ruin RabbitMQ
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in RabbitMQ goroutine: %v", r)
			}
		}()

		time.Sleep(5 * time.Second)
		log.Println("Starting RabbitMQ server...")
		rabbitServer := server.NewRabbitMQServer(database, config.RabbitMQURL, config.QueueName, config.ExchangeName)
		rabbitServer.Start()
	}()

	// run server
	serverDb := server.NewServer(database)
	serverDb.Start()
}
