package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IAuthService interface {
	CreateUser(user model.User) error
	GetUserByEmail(email string) (*model.User, error)
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

func (svc *authService) GetUserByEmail(email string) (*model.User, error) {
	return svc.userRepo.FindUserByEmail(email)
}
