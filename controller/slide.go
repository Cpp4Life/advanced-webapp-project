package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ISlideController interface {
	GetSlideById(c *gin.Context)
	CreateSlide(c *gin.Context)
	UpdateSlide(c *gin.Context)
	DeleteSlide(c *gin.Context)
}

type slideController struct {
	logger       *utils.Logger
	slideService service.ISlideService
}

func NewSlideController(logger *utils.Logger, slideSvc service.ISlideService) *slideController {
	return &slideController{
		logger:       logger,
		slideService: slideSvc,
	}
}

func (s *slideController) GetSlideById(c *gin.Context) {}

func (s *slideController) CreateSlide(c *gin.Context) {
	presId := c.Param("id")
	slideType := c.Query("type")

	s.logger.Warn(presId, slideType)

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		s.logger.Error(err.Error())
		return
	}

	var slide model.Slide
	json.Unmarshal(jsonData, &slide)

	presIdUint, _ := strconv.ParseUint(presId, 10, 64)
	slideTypeUint, _ := strconv.ParseUint(slideType, 10, 64)

	slide.PresentationId = uint(presIdUint)
	slide.Type = uint(slideTypeUint)

	slideId, err := s.slideService.CreateSlide(&slide)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "failed to create slide"})
		s.logger.Error(err.Error())
		return
	}

	content := slide.Content
	contentId, err := s.slideService.CreateContent(strconv.FormatInt(slideId, 10), content)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "failed to create content"})
		s.logger.Error(err.Error())
		return
	}

	options := content.Options
	_, err = s.slideService.CreateOption(strconv.FormatInt(contentId, 10), options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "failed to create options"})
		s.logger.Error(err.Error())
		return
	}

	slide.Id = uint(slideId)
	slide.Content.Id = uint(contentId)
	c.JSON(http.StatusCreated, map[string]any{
		"data": slide,
	})
}

func (s *slideController) UpdateSlide(c *gin.Context) {}

func (s *slideController) DeleteSlide(c *gin.Context) {}
