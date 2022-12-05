package repository

import (
	"advanced-webapp-project/model"
	"context"
	"database/sql"
)

type ISlideRepo interface {
	FindSlideById()
	InsertSlide(slide *model.Slide) (int64, error)
	InsertContent(slideId string, content *model.Content) (int64, error)
	InsertOption(contentId string, options []*model.Option) (int64, error)
	UpdateSlide() (int64, error)
	DeleteSlide() (int64, error)
}

type slideRepo struct {
	conn *sql.DB
}

func NewSlideRepo(sqldb *sql.DB) *slideRepo {
	return &slideRepo{
		conn: sqldb,
	}
}

func (db *slideRepo) FindSlideById() {}

func (db *slideRepo) InsertSlide(slide *model.Slide) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := db.conn.ExecContext(ctx, stmtInsertSlide, slide.PresentationId, slide.Type)
	if err != nil {
		return -1, nil
	}

	return result.LastInsertId()
}

func (db *slideRepo) InsertContent(slideId string, content *model.Content) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := db.conn.ExecContext(ctx, stmtInsertContent, slideId, content.Title, content.Meta)
	if err != nil {
		return -1, nil
	}

	return result.LastInsertId()
}

func (db *slideRepo) InsertOption(contentId string, options []*model.Option) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	for _, option := range options {
		_, err := db.conn.ExecContext(ctx, stmtInsertOption, option.Name, option.Image, contentId)
		if err != nil {
			return -1, nil
		}
	}

	return 0, nil
}

func (db *slideRepo) UpdateSlide() (int64, error) {
	return -1, nil
}

func (db *slideRepo) DeleteSlide() (int64, error) {
	return -1, nil
}
