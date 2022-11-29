package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IAuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	logger      *utils.Logger
	jwtService  service.IJWTService
	authService service.IAuthService
}

func NewAuthHandler(logger *utils.Logger, jwtSvc service.IJWTService, authSvc service.IAuthService) *authController {
	return &authController{
		logger:      logger,
		jwtService:  jwtSvc,
		authService: authSvc,
	}
}

func (ctl *authController) Login(c *gin.Context) {
	user := model.User{}
	// Bind JSON data to `User` model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		ctl.logger.Error(err.Error())
		return
	}

	userData, err := ctl.authService.VerifyCredential(user.Email, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "invalid credentials!"})
		ctl.logger.Error(err.Error())
		return
	}

	// Generate token
	userId := strconv.Itoa(int(userData.Id))
	generatedToken := ctl.jwtService.GenerateToken(userId, userData.Email)
	c.JSON(http.StatusOK, map[string]any{
		"user":  userData,
		"token": generatedToken,
	})

	return
}

func (ctl *authController) Register(c *gin.Context) {
	user := model.User{}
	// Bind JSON data to `User` model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		ctl.logger.Error(err.Error())
		return
	}

	// Check if email is whether taken or not
	isEmailCreated, _ := ctl.authService.GetUserByEmail(user.Email)
	if isEmailCreated != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": fmt.Sprintf("%s is in use", isEmailCreated.Email),
		})
		return
	}

	// Create user to db
	_, err := ctl.authService.CreateUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to create user!"})
		ctl.logger.Error(err.Error())
		return
	}

	ctl.logger.Warn(user)
	c.JSON(http.StatusCreated, map[string]any{
		"user": user,
	})

	return
}
