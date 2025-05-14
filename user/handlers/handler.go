package handlers

import usescases "user-auth/user/usecases"

type UserHttpHandler struct {
	Repo usescases.UserUsecase
}

type UserRabbitMQHandler struct {
	Repo usescases.UserEventUsecase
}

func NewUserHttpHandler(usecaseInput usescases.UserUsecase) *UserHttpHandler {
	return &UserHttpHandler{
		Repo: usecaseInput,
	}
}

func NewUserRabbitMQHandler(usecaseInput usescases.UserEventUsecase) *UserRabbitMQHandler {
	return &UserRabbitMQHandler{
		Repo: usecaseInput,
	}
}
