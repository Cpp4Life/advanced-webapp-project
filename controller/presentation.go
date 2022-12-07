package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IPresController interface {
	GetPresentationById(c *gin.Context)
	GetAllPresentations(c *gin.Context)
	CreatePresentation(c *gin.Context)
	UpdatePresentation(c *gin.Context)
	DeletePresentation(c *gin.Context)
}

type presController struct {
	logger      *utils.Logger
	jwtService  service.IJWTService
	presService service.IPresService
	userService service.IUserService
}

func NewPresController(logger *utils.Logger, jwtSvc service.IJWTService, presSvc service.IPresService, userSvc service.IUserService) *presController {
	return &presController{
		logger:      logger,
		jwtService:  jwtSvc,
		presService: presSvc,
		userService: userSvc,
	}
}

func (p *presController) GetPresentationById(c *gin.Context) {
	pres, err := p.presService.GetPresentationById(c.Param("id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "presentation not found!"})
		p.logger.Error(err.Error())
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, map[string]any{
		"data": pres,
	})

	return
}

func (p *presController) GetAllPresentations(c *gin.Context) {
	presList, err := p.presService.GetAllPresentations()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "no presentations found!"})
		p.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"data": presList,
	})

	return
}

func (p *presController) CreatePresentation(c *gin.Context) {
	userId := p.getUserId(c.GetHeader("Authorization"))
	pres := model.Pres{}
	if err := c.ShouldBindJSON(&pres); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		p.logger.Error(err.Error())
		return
	}

	randomId := utils.GenerateRandomNumber(8)
	presId, _ := strconv.ParseUint(randomId, 10, 64)
	pres.Id = uint(presId)
	err := p.presService.CreatePresentation(&pres, userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to create presentation"})
		p.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"data": pres,
	})

	return
}

func (p *presController) UpdatePresentation(c *gin.Context) {
	pres := model.Pres{}
	if err := c.ShouldBindJSON(&pres); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		p.logger.Error(err.Error())
		return
	}

	_, err := p.presService.UpdatePresentation(c.Param("id"), pres)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to update"})
		p.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "updated successfully",
	})
}

func (p *presController) DeletePresentation(c *gin.Context) {
	presId := c.Param("id")

	result, err := p.presService.DeletePresentation(presId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to delete presentation"})
		p.logger.Error(err.Error())
		return
	}

	if result == 0 {
		c.AbortWithStatusJSON(http.StatusConflict, map[string]any{"message": "presentation already deleted"})
		return
	}

	c.JSON(http.StatusNoContent, map[string]any{
		"message": "deleted successfully",
	})

	return
}

func (p *presController) getUserId(token string) string {
	claims, _ := p.jwtService.ExtractToken(token)
	return claims["user_id"].(string)
}
