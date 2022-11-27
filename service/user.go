package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IUserService interface {
	GetProfile(id string) (*model.User, error)
	UpdateProfile(id string, user model.User) (int64, error)
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

func (svc *userService) UpdateProfile(id string, user model.User) (int64, error) {
	return svc.userRepo.ModifyUserById(id, user)
}
