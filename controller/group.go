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
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "error while getting groups"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"groups": groups,
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
		"group": group,
	})

	return
}
