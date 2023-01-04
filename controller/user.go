package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type IUserController interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	ChangePassword(c *gin.Context)
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

func (u *userController) ChangePassword(c *gin.Context) {
	type req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	jsonData := req{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		u.logger.Error(err.Error())
		return
	}

	if jsonData.OldPassword == "" || jsonData.NewPassword == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "old password and new password must be filled in"})
		return
	}

	if jsonData.OldPassword == jsonData.NewPassword {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "old password must not be the same as new password"})
		return
	}

	userId := u.getUserId(c.GetHeader("Authorization"))
	userProfile, _ := u.userService.GetProfile(userId)
	if err := bcrypt.CompareHashAndPassword([]byte(userProfile.SavedPassword), []byte(jsonData.OldPassword)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "old passwords do not match"})
		return
	}

	if userProfile.IsSocial {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "cannot change password due to email is logged in via social account"})
		return
	}

	_, err := u.userService.ChangePassword(userId, jsonData.NewPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to change password"})
		u.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{"message": "updated successfully"})
	return
}

func (u *userController) getUserId(token string) string {
	claims, _ := u.jwtService.ExtractToken(token)
	return claims["user_id"].(string)
}
