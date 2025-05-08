package entities

import (
	"errors"
	"log/slog"
	"time"
	"user-auth/utils"

	"github.com/google/uuid"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	Role      string `json:"role" gorm:"default:'user'"` // default role is user
	AuthKey   string `json:"auth_key"`
	Avatar    string `json:"avatar"` //It will be saved in the S3 bucket

}

func (u *User) BeforeCreate() error {
	u.ID = uuid.New().String()
	u.AuthKey = uuid.New().String()
	u.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	u.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		slog.Error("Failed to hash password", err)
		return errors.New("failed to hash password")
	}
	u.Password = hashedPassword

	return nil
}
