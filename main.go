package main

import (
	"log"
	"user-auth/db"
	"user-auth/server"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	server := server.NewServer(database)
	server.Start()
}
