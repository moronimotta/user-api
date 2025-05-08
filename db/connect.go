// db/connect.go
package db

import (
	"log"
	user "user-auth/user/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (Database, error) {
	config := Config{}
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(config.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	return &GormDatabase{DB: db}, nil
}
