package repositories

import "user-auth/user/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByID(id string) (*entities.User, error)
	GetAllUsers() ([]entities.User, error)
	UpdateUser(user *entities.User) error
	DeleteUser(id string) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByEmailAndPassword(email, password string) (*entities.User, error)
}
