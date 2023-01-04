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
	ForgotPassword(c *gin.Context)
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
		if isEmailCreated.IsSocial && user.IsSocial {
			userId := strconv.Itoa(int(isEmailCreated.Id))
			generatedToken := ctl.jwtService.GenerateToken(userId, isEmailCreated.Email)
			c.JSON(http.StatusOK, map[string]any{
				"user":  isEmailCreated,
				"token": generatedToken,
			})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
				"message": fmt.Sprintf("%s is in use", isEmailCreated.Email),
			})
			return
		}
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

	if user.IsSocial {
		userId := strconv.Itoa(int(user.Id))
		generatedToken := ctl.jwtService.GenerateToken(userId, user.Email)
		c.JSON(http.StatusOK, map[string]any{
			"user":  user,
			"token": generatedToken,
		})
	} else {
		email := service.Message{
			URL:       "https://kameyoko.up.railway.app/api/v1/auth/verify-email/" + code,
			FullName:  user.FullName,
			Subject:   "Your account verification code",
			Paragraph: "Please verify your account to be able to login",
		}

		go func(user *model.User, email *service.Message) {
			ctl.mailService.SendMail(user, email)
		}(&user, &email)

		c.JSON(http.StatusCreated, map[string]any{
			"message": "We sent an email with a verification code to " + user.Email,
		})
	}
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

func (ctl *authController) ForgotPassword(c *gin.Context) {
	type req struct {
		Email string `json:"email"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		ctl.logger.Error(err.Error())
		return
	}

	user, _ := ctl.authService.GetUserByEmail(body.Email)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"message": "email not found"})
		return
	}

	user.IsVerified, _ = ctl.authService.GetVerifiedStatusByEmail(body.Email)
	if !user.IsVerified && !user.IsSocial {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "verify your account first in order to reset password"})
		return
	}

	if user.IsSocial {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "cannot reset password due to email is logged in via social account"})
		return
	}

	newPassword := randstr.String(10)
	_, err := ctl.authService.UpdatePassword(strconv.Itoa(int(user.Id)), newPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to reset password"})
		ctl.logger.Error(err.Error())
		return
	}

	email := service.Message{
		URL:       "http://localhost:3000/login",
		FullName:  user.FullName,
		Subject:   "New password for account",
		Paragraph: fmt.Sprintf("Here is your new password %s. Please keep it secret!", newPassword),
	}

	go func() {
		ctl.mailService.SendMail(user, &email)
	}()

	c.JSON(http.StatusOK, map[string]any{
		"message": "new password has been sent to " + user.Email,
	})
}
