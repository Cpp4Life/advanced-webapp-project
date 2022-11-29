package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IUserController interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

type userController struct {
	logger       *utils.Logger
	jwtService   service.IJWTService
	userService  service.IUserService
	groupService service.IGroupService
}

func NewUserController(logger *utils.Logger, jwtSvc service.IJWTService, userSvc service.IUserService, groupSvc service.IGroupService) *userController {
	return &userController{
		logger:       logger,
		jwtService:   jwtSvc,
		userService:  userSvc,
		groupService: groupSvc,
	}
}

// godoc
// @Security Token
// @Summary Get profile
// @Description Get user profile after login
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 404
// @Router /accounts/profile [get]
func (u *userController) GetProfile(c *gin.Context) {
	userId := u.getUserId(c.GetHeader("Authorization"))
	user, err := u.userService.GetProfile(userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"message": "profile not found!"})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"user": user,
	})

	return
}

func (u *userController) UpdateProfile(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		u.logger.Error(err.Error())
		return
	}

	userId := u.getUserId(c.GetHeader("Authorization"))
	_, err := u.userService.UpdateProfile(userId, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to update user"})
		u.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusNoContent, map[string]any{
		"message": "Update successfully",
	})
	return
}

func (u *userController) getUserId(token string) string {
	claims, _ := u.jwtService.ExtractToken(token)
	return claims["user_id"].(string)
}
