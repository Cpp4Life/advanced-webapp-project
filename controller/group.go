package controller

import (
	"advanced-webapp-project/helper"
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IGroupController interface {
	GetAllGroups(c *gin.Context)
	CreateGroup(c *gin.Context)
	GetCreatedGroupsByUserId(c *gin.Context)
	GetJoinedGroupsByUserId(c *gin.Context)
	GetGroupMemberDetailsByGroupId(c *gin.Context)
}

type groupController struct {
	logger       *helper.Logger
	jwtService   service.IJWTService
	groupService service.IGroupService
	userService  service.IUserService
}

func NewGroupController(logger *helper.Logger, jwtSvc service.IJWTService, groupSvc service.IGroupService, userSvc service.IUserService) *groupController {
	return &groupController{
		logger:       logger,
		jwtService:   jwtSvc,
		groupService: groupSvc,
		userService:  userSvc,
	}
}

func (g *groupController) GetAllGroups(c *gin.Context) {
	groups, err := g.groupService.GetAllGroups()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "no groups found!"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"groups_data": groups,
	})

	return
}

func (g *groupController) CreateGroup(c *gin.Context) {
	group := model.Group{}
	if err := c.ShouldBindJSON(&group); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		g.logger.Error(err.Error())
		return
	}

	claims, _ := g.jwtService.ExtractToken(c.GetHeader("Authorization"))
	userId := claims["user_id"].(string)
	userData, err := g.userService.GetProfile(userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "please login first!"})
		return
	}

	_, err = g.groupService.CreateGroup(&group, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to create group!"})
		g.logger.Error(err.Error())
		return
	}

	group.Owner = userData
	c.JSON(http.StatusCreated, map[string]any{
		"group_data": group,
	})

	return
}

func (g *groupController) GetCreatedGroupsByUserId(c *gin.Context) {
	userId := g.getUserId(c.GetHeader("Authorization"))
	groups, err := g.groupService.GetCreatedGroupsByUserId(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "no groups found!"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"groups_data": groups,
	})

	return
}

func (g *groupController) GetJoinedGroupsByUserId(c *gin.Context) {
	userId := g.getUserId(c.GetHeader("Authorization"))
	groups, err := g.groupService.GetJoinedGroupsByUserId(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "no groups found!"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"groups_data": groups,
	})

	return
}

func (g *groupController) GetGroupMemberDetailsByGroupId(c *gin.Context) {
	groups, err := g.groupService.GetGroupMemberDetailsByGroupId(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "no groups found!"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"groups_data": groups,
	})

	return
}

func (g *groupController) getUserId(token string) string {
	claims, _ := g.jwtService.ExtractToken(token)
	return claims["user_id"].(string)
}
