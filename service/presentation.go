package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IPresService interface {
	GetPresentationById(presId string) (*model.Pres, error)
	GetAllPresentations() ([]*model.Pres, error)
	CreatePresentation(pres *model.Pres, userId string) error
	UpdatePresentation(presId string, data model.Pres) (int64, error)
	DeletePresentation(presId string) (int64, error)
	PresentGroup(data model.GroupPresInfo) (int64, error)
}

type presService struct {
	presRepo repository.IPresRepo
}

func NewPresService(presRepo repository.IPresRepo) *presService {
	return &presService{
		presRepo: presRepo,
	}
}

func (svc *presService) GetPresentationById(presId string) (*model.Pres, error) {
	return svc.presRepo.FindPresentationById(presId)
}

func (svc *presService) GetAllPresentations() ([]*model.Pres, error) {
	return svc.presRepo.FindAllPresentations()
}

func (svc *presService) CreatePresentation(pres *model.Pres, userId string) error {
	return svc.presRepo.InsertPresentation(pres, userId)
}

func (svc *presService) UpdatePresentation(presId string, data model.Pres) (int64, error) {
	return svc.presRepo.UpdatePresentation(presId, data)
}

func (svc *presService) DeletePresentation(presId string) (int64, error) {
	return svc.presRepo.DeletePresentation(presId)
}

func (svc *presService) PresentGroup(data model.GroupPresInfo) (int64, error) {
	return svc.presRepo.PresentGroup(data)
}
