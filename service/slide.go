package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type ISlideService interface {
	GetAllSlides(presId string)
	CreateSlide(slide *model.Slide) (int64, error)
	CreateContent(slideId string, content *model.Content) (int64, error)
	CreateOption(contentId string, options []*model.Option) (int64, error)
	UpdateSlide(presId string, slide model.Slide) (int64, error)
	UpdateContent(slideId string, content model.Content) (int64, error)
	UpdateOptions(contentId string, options []*model.Option) (int64, error)
	DeleteSlide() (int64, error)
}

type slideService struct {
	slideRepo repository.ISlideRepo
}

func NewSlideService(slideRepo repository.ISlideRepo) *slideService {
	return &slideService{
		slideRepo: slideRepo,
	}
}

func (svc *slideService) GetAllSlides(presId string) {

}

func (svc *slideService) CreateSlide(slide *model.Slide) (int64, error) {
	return svc.slideRepo.InsertSlide(slide)
}

func (svc *slideService) CreateContent(slideId string, content *model.Content) (int64, error) {
	return svc.slideRepo.InsertContent(slideId, content)
}

func (svc *slideService) CreateOption(contentId string, options []*model.Option) (int64, error) {
	return svc.slideRepo.InsertOption(contentId, options)
}

func (svc *slideService) UpdateSlide(presId string, slide model.Slide) (int64, error) {
	return svc.slideRepo.UpdateSlide(presId, slide)
}

func (svc *slideService) UpdateContent(slideId string, content model.Content) (int64, error) {
	return svc.slideRepo.UpdateContent(slideId, content)
}

func (svc *slideService) UpdateOptions(contentId string, options []*model.Option) (int64, error) {
	return svc.slideRepo.UpdateOptions(contentId, options)
}

func (svc *slideService) DeleteSlide() (int64, error) {
	return svc.slideRepo.DeleteSlide()
}
