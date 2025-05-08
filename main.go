package main

import (
	"log"
	"log/slog"
	"user-auth/db"
	"user-auth/server"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	slog.Info("Connected to DB successfully")
	server := server.NewServer(database)
	server.Start()
}
