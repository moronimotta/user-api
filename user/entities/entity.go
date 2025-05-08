package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	Role      string `json:"role"`
	Avatar    string `json:"avatar"` //It will be saved in the S3 bucket
}

func (u *User) BeforeCreate() error {
	u.ID = uuid.New().String()
	u.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	u.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

// create Dummy user for testing purposes
var UserEntity = User{
	Name:     "John Doe",
	Email:    "example@example.com",
	Password: "password",
	Role:     "user",
	Avatar:   "https://example.com/avatar.jpg",
}
