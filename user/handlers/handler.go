package handlers

import usescases "user-auth/user/usecases"

type userHttpHandler struct {
	Repo usescases.UserUsecase
}

func NewUserHttpHandler(usecaseInput usescases.UserUsecase) *userHttpHandler {
	return &userHttpHandler{
		Repo: usecaseInput,
	}
}
