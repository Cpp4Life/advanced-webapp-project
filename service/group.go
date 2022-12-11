package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type IGroupService interface {
	GetAllGroups() ([]*model.Group, error)
	CreateGroup(group *model.Group, userId string) (int64, error)
	GetCreatedGroupsByUserId(userId string) ([]*model.Group, error)
	GetJoinedGroupsByUserId(userId string) ([]*model.GroupUser, error)
	GetGroupMemberDetailsByGroupId(groupId string) ([]*model.GroupUser, error)
	GetGroupById(groupId string) (*model.Group, error)
	UpdateUserRole(groupId, userId, role string) (int64, error)
	GetUserRole(groupId, userId string) (string, error)
	AddMemberToGroup(groupId string, member model.Member) (int64, error)
	DeleteMember(groupId, userId string) (int64, error)
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

func (svc *groupService) GetCreatedGroupsByUserId(userId string) ([]*model.Group, error) {
	return svc.groupRepo.FindCreatedGroupsByUserId(userId)
}

func (svc *groupService) GetJoinedGroupsByUserId(userId string) ([]*model.GroupUser, error) {
	return svc.groupRepo.FindJoinedGroupsByUserId(userId)
}

func (svc *groupService) GetGroupMemberDetailsByGroupId(groupId string) ([]*model.GroupUser, error) {
	return svc.groupRepo.FindGroupMemberDetailsByGroupId(groupId)
}

func (svc *groupService) GetGroupById(groupId string) (*model.Group, error) {
	return svc.groupRepo.FindGroupById(groupId)
}

func (svc *groupService) UpdateUserRole(groupId, userId, role string) (int64, error) {
	return svc.groupRepo.UpdateUserRole(groupId, userId, role)
}

func (svc *groupService) GetUserRole(groupId, userId string) (string, error) {
	return svc.groupRepo.FindUserRole(groupId, userId)
}

func (svc *groupService) AddMemberToGroup(groupId string, member model.Member) (int64, error) {
	return svc.groupRepo.InsertMemberToGroup(groupId, member)
}

func (svc *groupService) DeleteMember(groupId, userId string) (int64, error) {
	return svc.groupRepo.DeleteMember(groupId, userId)
}
