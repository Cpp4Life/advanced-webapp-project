package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IUserService interface {
	GetProfile(id string) (*model.User, error)
}

type userService struct {
	userRepo repository.IUserRepo
}

func NewUserService(userRepo repository.IUserRepo) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (svc *userService) GetProfile(id string) (*model.User, error) {
	return svc.userRepo.FindUserById(id)
}
