package handlers

import (
	usescases "user-auth/user/usecases"
)

type UserHttpHandler struct {
	Repo usescases.UserUsecase
}

func NewUserHttpHandler(usecaseInput usescases.UserUsecase) *UserHttpHandler {
	return &UserHttpHandler{
		Repo: usecaseInput,
	}
}
