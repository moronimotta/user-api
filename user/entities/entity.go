package entities

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"
	"user-auth/utils"

	"github.com/google/uuid"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"password"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
	Role       string `json:"role" gorm:"default:'customer'"`
	AuthKey    string `json:"auth_key"`
	Avatar     string `json:"avatar"`      //It will be saved in the S3 bucket
	ExternalID string `json:"external_id"` // Link to Finance API
}

const (
	RoleAdmin    = "admin"
	RoleBarber   = "barber"
	RoleCustomer = "customer"
)

// CanonicalRole keeps API and token role claims within the supported role set.
// The historical "user" role is treated as a customer for existing accounts.
func CanonicalRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case RoleAdmin:
		return RoleAdmin
	case RoleBarber:
		return RoleBarber
	default:
		return RoleCustomer
	}
}

// UserResponse is the public profile shape. It intentionally excludes password
// hashes, auth keys, and payment-provider identifiers.
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewUserResponse(user *User) UserResponse {
	return UserResponse{
		ID: user.ID, Name: user.Name, Email: user.Email,
		Role: CanonicalRole(user.Role), Avatar: user.Avatar,
		CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	}
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

func (u *User) Unmarshal(data []byte) error {
	if !json.Valid(data) {
		slog.Error("Invalid JSON data")
		return errors.New("invalid JSON data")
	}
	if err := json.Unmarshal(data, u); err != nil {
		slog.Error("Failed to unmarshal user data", err)
		return errors.New("failed to unmarshal user data")
	}
	return nil
}
