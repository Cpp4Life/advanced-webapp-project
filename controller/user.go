package controller

import (
	"advanced-webapp-project/helper"
	"advanced-webapp-project/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IUserController interface {
	GetProfile(c *gin.Context)
}

type userController struct {
	logger      *helper.Logger
	jwtService  service.IJWTService
	userService service.IUserService
}

func NewUserController(logger *helper.Logger, jwtSvc service.IJWTService, userSvc service.IUserService) *userController {
	return &userController{
		logger:      logger,
		jwtService:  jwtSvc,
		userService: userSvc,
	}
}

func (u *userController) GetProfile(c *gin.Context) {
	claims, _ := u.jwtService.ExtractToken(c.GetHeader("Authorization"))
	userId := claims["user_id"].(string)
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
