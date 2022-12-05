package service

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/repository"
)

type ISlideService interface {
	GetSlideById()
	CreateSlide(slide *model.Slide) (int64, error)
	CreateContent(slideId string, content *model.Content) (int64, error)
	CreateOption(contentId string, options []*model.Option) (int64, error)
	UpdateSlide() (int64, error)
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

func (svc *slideService) GetSlideById() {

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

func (svc *slideService) UpdateSlide() (int64, error) {
	return svc.slideRepo.UpdateSlide()
}

func (svc *slideService) DeleteSlide() (int64, error) {
	return svc.slideRepo.DeleteSlide()
}
