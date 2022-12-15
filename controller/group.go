package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"net/http"
	"strconv"
)

type IGroupController interface {
	GetAllGroups(c *gin.Context)
	CreateGroup(c *gin.Context)
	GetCreatedGroupsByUserId(c *gin.Context)
	GetJoinedGroupsByUserId(c *gin.Context)
	GetGroupMemberDetailsByGroupId(c *gin.Context)
	GetGroupById(c *gin.Context)
	UpdateUserRole(c *gin.Context)
	AddMemberToGroup(c *gin.Context)
	DeleteMember(c *gin.Context)
	InviteMember(c *gin.Context)
}

type groupController struct {
	logger       *utils.Logger
	jwtService   service.IJWTService
	groupService service.IGroupService
	userService  service.IUserService
	authService  service.IAuthService
	mailService  service.IMailerService
}

func NewGroupController(logger *utils.Logger, jwtSvc service.IJWTService, groupSvc service.IGroupService, userSvc service.IUserService, authSvc service.IAuthService, mailSvc service.IMailerService) *groupController {
	return &groupController{
		logger:       logger,
		jwtService:   jwtSvc,
		groupService: groupSvc,
		userService:  userSvc,
		authService:  authSvc,
		mailService:  mailSvc,
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

	userId := g.getUserId(c.GetHeader("Authorization"))
	userData, err := g.userService.GetProfile(userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "please login first!"})
		return
	}

	group.Link = randstr.String(32)
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

func (g *groupController) GetGroupById(c *gin.Context) {
	groupId := c.Param("id")
	group, err := g.groupService.GetGroupById(groupId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"message": "no group found!"})
		g.logger.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"group_data": group,
	})

	return
}

func (g *groupController) UpdateUserRole(c *gin.Context) {
	groupId := c.Param("id")
	userId := c.Query("userId")
	role := c.Query("role")

	if _, err := g.userService.GetProfile(userId); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"message": "user not found"})
		g.logger.Error(err.Error())
		return
	}

	if _, err := g.groupService.GetGroupById(groupId); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"message": "group not found"})
		g.logger.Error(err.Error())
		return
	}

	if _, err := g.groupService.UpdateUserRole(groupId, userId, role); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to assign role"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "successfully updated",
	})

	return
}

func (g *groupController) AddMemberToGroup(c *gin.Context) {
	var member model.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		g.logger.Error(err.Error())
		return
	}

	groupId := c.Param("id")
	user, _ := g.authService.GetUserByEmail(member.Email)
	member.UserId = user.Id

	_, err := g.groupService.AddMemberToGroup(groupId, member)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "user already a member"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"message": "invited successfully",
	})

	return
}

func (g *groupController) DeleteMember(c *gin.Context) {
	var member model.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		g.logger.Error(err.Error())
		return
	}

	groupId := c.Param("id")
	userId := g.getUserId(c.GetHeader("Authorization"))

	userRole, err := g.groupService.GetUserRole(groupId, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "user with undefined role"})
		g.logger.Error(err.Error())
		return
	}

	if !g.isOwner(userRole) && !g.isCoOwner(userRole) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "user doesn't have permission to invite member"})
		return
	}

	_, err = g.groupService.DeleteMember(groupId, strconv.Itoa(int(member.UserId)))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to delete user"})
		g.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "deleted successfully",
	})

	return
}

func (g *groupController) InviteMember(c *gin.Context) {
	type invitationReq struct {
		Email string `json:"email"`
		Link  string `json:"link"`
	}

	var data invitationReq
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		g.logger.Error(err.Error())
		return
	}

	groupId := c.Param("id")
	userId := g.getUserId(c.GetHeader("Authorization"))
	userRole, err := g.groupService.GetUserRole(groupId, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "user with undefined role"})
		g.logger.Error(err.Error())
		return
	}

	if !g.isOwner(userRole) && !g.isCoOwner(userRole) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "user doesn't have permission to invite member"})
		return
	}

	email := service.Message{
		URL:       data.Link,
		FullName:  "there",
		Subject:   "Invitation link to group",
		Paragraph: "Please click the link below to join our group",
	}

	go func(user *model.User, email *service.Message) {
		g.mailService.SendMail(user, email)
	}(&model.User{Email: data.Email}, &email)

	c.JSON(http.StatusCreated, map[string]any{
		"message": "An invitation has been sent to " + data.Email,
	})

	return
}

func (g *groupController) isOwner(role string) bool {
	return role == "1"
}

func (g *groupController) isCoOwner(role string) bool {
	return role == "2"
}

func (g *groupController) getUserId(token string) string {
	claims, _ := g.jwtService.ExtractToken(token)
	return claims["user_id"].(string)
}
