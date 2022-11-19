package controller

import (
	"advanced-webapp-project/helper"
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	logger     helper.ILogger
	jwtService service.IJWTService
}

func NewAuthHandler(log helper.ILogger, jwtSvc service.IJWTService) *authController {
	return &authController{
		logger:     log,
		jwtService: jwtSvc,
	}
}

func (ctl *authController) Login(c *gin.Context) {

}

func (ctl *authController) Register(c *gin.Context) {
	user := model.User{}
	validate := validator.New()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		ctl.logger.Error(err.Error())
		return
	}

	ctl.logger.Info(user)

	if err := validate.Struct(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		ctl.logger.Error(err.Error())
		return
	}

	return
}
