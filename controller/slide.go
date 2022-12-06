package controller

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ISlideController interface {
	GetAllSlides(c *gin.Context)
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

func (s *slideController) GetAllSlides(c *gin.Context) {
	presId := c.Param("id")

	slides, err := s.slideService.GetAllSlides(presId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "slides not found!"})
		s.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"slides": slides,
	})
}

func (s *slideController) CreateSlide(c *gin.Context) {
	presId := c.Param("id")

	var slide model.Slide
	if err := c.ShouldBindJSON(&slide); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		s.logger.Error(err.Error())
		return
	}

	presIdUint, _ := strconv.ParseUint(presId, 10, 64)
	slide.PresentationId = uint(presIdUint)

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

func (s *slideController) UpdateSlide(c *gin.Context) {
	presId := c.Param("id")
	slideId := c.Param("slide_id")
	contentId := c.Query("content_id")

	var slide model.Slide
	if err := c.ShouldBindJSON(&slide); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		s.logger.Error(err.Error())
		return
	}

	slideIdUint, _ := strconv.ParseUint(slideId, 10, 64)
	slide.Id = uint(slideIdUint)

	_, err := s.slideService.UpdateSlide(presId, slide)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to update slide"})
		s.logger.Error(err.Error())
		return
	}

	_, err = s.slideService.UpdateContent(slideId, *slide.Content)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to update content"})
		s.logger.Error(err.Error())
		return
	}

	_, err = s.slideService.UpdateOptions(contentId, slide.Content.Options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to update options"})
		s.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "updated successfully",
	})
}

func (s *slideController) DeleteSlide(c *gin.Context) {}
