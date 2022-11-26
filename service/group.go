package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IGroupService interface {
	GetAllGroups() ([]*model.Group, error)
	CreateGroup(group *model.Group, userId string) (int64, error)
}

type groupService struct {
	groupRepo repository.IGroupRepo
}

func NewGroupService(groupRepo repository.IGroupRepo) *groupService {
	return &groupService{
		groupRepo: groupRepo,
	}
}

func (svc *groupService) GetAllGroups() ([]*model.Group, error) {
	return svc.groupRepo.FindAll()
}

func (svc *groupService) CreateGroup(group *model.Group, userId string) (int64, error) {
	return svc.groupRepo.InsertGroup(group, userId)
}
