package user

import "github.com/google/uuid"

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
	return nil
}
