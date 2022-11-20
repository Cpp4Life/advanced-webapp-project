package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IAuthService interface {
	CreateUser(user model.User) error
}

type authService struct {
	userRepo repository.IUserRepo
}

func NewAuthService(userRepo repository.IUserRepo) *authService {
	return &authService{
		userRepo: userRepo,
	}
}

func (svc *authService) CreateUser(user model.User) error {
	return svc.userRepo.InsertUser(user)
}
