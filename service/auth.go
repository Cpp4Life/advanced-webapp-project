package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IAuthService interface {
	CreateUser(user *model.User) (int64, error)
	GetUserByEmail(email string) (*model.User, error)
	VerifyCredential(email, password string) (*model.User, error)
	GetVerifiedStatusByEmail(email string) (bool, error)
	UpdateVerifiedStatus(verificationCode string) (int64, error)
	UpdatePassword(id, password string) (int64, error)
}

type authService struct {
	userRepo repository.IUserRepo
}

func NewAuthService(userRepo repository.IUserRepo) *authService {
	return &authService{
		userRepo: userRepo,
	}
}

func (svc *authService) CreateUser(user *model.User) (int64, error) {
	return svc.userRepo.InsertUser(user)
}

func (svc *authService) GetUserByEmail(email string) (*model.User, error) {
	return svc.userRepo.FindUserByEmail(email)
}

func (svc *authService) VerifyCredential(email, password string) (*model.User, error) {
	return svc.userRepo.VerifyCredential(email, password)
}

func (svc *authService) GetVerifiedStatusByEmail(email string) (bool, error) {
	return svc.userRepo.FindVerifiedStatusByEmail(email)
}

func (svc *authService) UpdateVerifiedStatus(verificationCode string) (int64, error) {
	return svc.userRepo.UpdateVerifiedStatus(verificationCode)
}

func (svc *authService) UpdatePassword(id, password string) (int64, error) {
	return svc.userRepo.UpdatePassword(id, password)
}
