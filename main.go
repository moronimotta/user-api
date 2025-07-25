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
	err := confs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// connect to database Postgres
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// initialize Redis
	redisServer := server.NewRedisServer(database)
	redisClient := redisServer.Start()
	defer redisServer.Close()

	// run RabbitMQ
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in RabbitMQ goroutine: %v", r)
			}
		}()

		time.Sleep(5 * time.Second)
		log.Println("Starting RabbitMQ server...")
		rabbitServer := server.NewRabbitMQServer(database, redisClient)
		rabbitServer.Start()
	}()

	// run server with Redis client
	serverDb := server.NewServer(database, redisClient)
	serverDb.Start()
}
