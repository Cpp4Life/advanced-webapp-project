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
	UpdateOptionVote(c *gin.Context)
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
	slideId := utils.GenerateRandomNumber(8)
	contentId := utils.GenerateRandomNumber(8)

	var slide model.Slide
	if err := c.ShouldBindJSON(&slide); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		s.logger.Error(err.Error())
		return
	}

	slide.PresentationId = utils.Str2Uint(presId)
	slide.Id = utils.Str2Uint(slideId)

	err := s.slideService.CreateSlide(&slide)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "failed to create slide"})
		s.logger.Error(err.Error())
		return
	}

	content := slide.Content
	content.Id = utils.Str2Uint(contentId)
	err = s.slideService.CreateContent(slideId, content)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "failed to create content"})
		s.logger.Error(err.Error())
		return
	}

	switch slide.Type {
	case 1:
		options := content.Options
		err = s.slideService.CreateOption(contentId, options)
	case 8:
		heading := content.Heading
		err = s.slideService.CreateHeading(contentId, heading)
	case 9:
		paragraph := content.Paragraph
		err = s.slideService.CreateParagraph(contentId, paragraph)
	default:
		s.logger.Warn("Create slide was interrupted")
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": "failed to create slide content"})
		s.logger.Error(err.Error())
		return
	}

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

	switch slide.Type {
	case 1:
		_, err = s.slideService.UpdateOptions(contentId, slide.Content.Options)
	case 8:
		_, err = s.slideService.UpdateHeading(contentId, slide.Content.Heading)
	case 9:
		_, err = s.slideService.UpdateParagraph(contentId, slide.Content.Paragraph)
	default:
		s.logger.Warn("Update slide was interrupted")
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to update sub-content"})
		s.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "updated successfully",
	})
}

func (s *slideController) UpdateOptionVote(c *gin.Context) {
	contentId := c.Param("content_id")
	optionId := c.Query("option_id")

	res, err := s.slideService.UpdateOptionVote(contentId, optionId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to vote option"})
		s.logger.Error(err.Error())
		return
	}

	if res == 0 {
		c.AbortWithStatusJSON(http.StatusConflict, map[string]any{"message": "failed to vote option"})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "voted successfully",
	})
}

func (s *slideController) DeleteSlide(c *gin.Context) {
	presId := c.Param("id")
	slideId := c.Param("slide_id")

	res, err := s.slideService.DeleteSlide(presId, slideId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"message": "failed to delete slide"})
		s.logger.Error(err.Error())
		return
	}

	if res == 0 {
		c.AbortWithStatusJSON(http.StatusConflict, map[string]any{"message": "slide already deleted"})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "deleted successfully",
	})
}
