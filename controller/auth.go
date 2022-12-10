package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"net/http"
	"strconv"
)

type IAuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	VerifyEmail(c *gin.Context)
}

type authController struct {
	logger      *utils.Logger
	jwtService  service.IJWTService
	authService service.IAuthService
	mailService service.IMailerService
}

func NewAuthHandler(logger *utils.Logger, jwtSvc service.IJWTService, authSvc service.IAuthService, mailSvc service.IMailerService) *authController {
	return &authController{
		logger:      logger,
		jwtService:  jwtSvc,
		authService: authSvc,
		mailService: mailSvc,
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "account not found!"})
		ctl.logger.Error(err.Error())
		return
	}

	isVerified, err := ctl.authService.GetVerifiedStatusByEmail(user.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "account not found!"})
		ctl.logger.Error(err.Error())
		return
	}

	if !isVerified {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "please verify your email first!"})
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

	code := randstr.String(20)
	verificationCode := utils.Encode(code)
	user.VerificationCode = verificationCode

	// Create user to db
	_, err := ctl.authService.CreateUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to create user!"})
		ctl.logger.Error(err.Error())
		return
	}

	email := service.Message{
		URL:      "http://localhost:7777/api/v1/auth/verify-email/" + code,
		FullName: user.FullName,
		Subject:  "Your account verification code",
	}

	go func(user *model.User, email *service.Message) {
		ctl.mailService.SendMail(user, email)
	}(&user, &email)

	c.JSON(http.StatusCreated, map[string]any{
		"message": "We sent an email with a verification code to " + user.Email,
	})

	return
}

func (ctl *authController) VerifyEmail(c *gin.Context) {
	code := c.Param("code")
	verificationCode := utils.Encode(code)

	result, err := ctl.authService.UpdateVerifiedStatus(verificationCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to verify account!"})
		ctl.logger.Error(err.Error())
		return
	}

	if result == 0 {
		c.AbortWithStatusJSON(http.StatusConflict, map[string]any{"message": "account already verified"})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/login?redirect=")
}
