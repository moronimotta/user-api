package usecases

import (
	"errors"
	"user-auth/user/entities"
	"user-auth/user/repositories"

	messageWorker "github.com/moronimotta/message-worker-module"
)

type UserEventUsecase struct {
	UserRepository repositories.UserRepository
}

func NewUserEventUsecase(userRepository repositories.UserRepository) *UserEventUsecase {
	return &UserEventUsecase{
		UserRepository: userRepository,
	}
}

func (u *UserEventUsecase) EventBus(event string) error {

	var eventData messageWorker.Event
	if err := eventData.Unmarshal([]byte(event)); err != nil {
		return err
	}
	userEntity := &entities.User{}

	dataBytes, ok := eventData.Data.([]byte)
	if !ok {
		return errors.New("event data is not of type []byte")
	}
	if err := userEntity.Unmarshal(dataBytes); err != nil {
		return err
	}
	switch eventData.Event {
	case "user.updated":
		u.UserRepository.UpdateUser(userEntity)
	default:
		return errors.New("event not found")
	}
	return nil
}
