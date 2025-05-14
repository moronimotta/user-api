// db/connect.go
package db

import (
	"log"
	"user-auth/confs"
	user "user-auth/user/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(confs confs.Config) (Database, error) {

	db, err := gorm.Open(postgres.Open(confs.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	return &GormDatabase{DB: db}, nil
}
