package usecases

import (
	"user-auth/user/entities"
	"user-auth/user/repositories"
)

type UserUsecase struct {
	UserRepository repositories.UserRepository
}

func NewUserUsecase(userRepository repositories.UserRepository) *UserUsecase {
	return &UserUsecase{
		UserRepository: userRepository,
	}
}

func (u *UserUsecase) GetUserByID(id string) (*entities.User, error) {
	user, err := u.UserRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) CreateUser(user *entities.User) error {
	err := user.BeforeCreate()
	if err != nil {
		return err
	}
	err = u.UserRepository.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) GetAllUsers() ([]entities.User, error) {
	users, err := u.UserRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserUsecase) UpdateUser(user *entities.User) error {
	err := u.UserRepository.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) DeleteUser(id string) error {
	err := u.UserRepository.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) GetUserByEmail(email string) (*entities.User, error) {
	user, err := u.UserRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) GetUserByEmailAndPassword(email, password string) (*entities.User, error) {
	user, err := u.UserRepository.GetUserByEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
