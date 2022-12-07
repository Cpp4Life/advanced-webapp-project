package repository

import (
	"advanced-webapp-project/model"
	"context"
	"database/sql"
)

type ISlideRepo interface {
	FindAllSlides(presId string) ([]*model.Slide, error)
	InsertSlide(slide *model.Slide) error
	InsertContent(slideId string, content *model.Content) error
	InsertOption(contentId string, options []*model.Option) error
	UpdateSlide(presId string, slide model.Slide) (int64, error)
	UpdateContent(slideId string, content model.Content) (int64, error)
	UpdateOptions(contentId string, options []*model.Option) (int64, error)
	UpdateOptionVote(contentId string, optionId string) (int64, error)
	DeleteSlide(presId, slideId string) (int64, error)
}

type slideRepo struct {
	conn *sql.DB
}

func NewSlideRepo(sqldb *sql.DB) *slideRepo {
	return &slideRepo{
		conn: sqldb,
	}
}

func (db *slideRepo) FindAllSlides(presId string) ([]*model.Slide, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx, stmtSelectAllSlides, presId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slides []*model.Slide
	for rows.Next() {
		var slide model.Slide
		var content model.Content
		var option model.Option
		var index = len(slides) - 1
		if err = rows.Scan(
			&slide.Id,
			&slide.Type,
			&content.Id,
			&content.Title,
			&content.Meta,
			&option.Id,
			&option.Name,
			&option.Image,
			&option.TotalVotes,
		); err != nil {
			return nil, err
		}

		content.Options = append(content.Options, &option)
		slide.Content = &content

		if len(slides) == 0 {
			slides = append(slides, &slide)
		} else if slides[index].Id == slide.Id {
			slides[index].Content.Options = append(slides[index].Content.Options, &option)
		} else {
			slides = append(slides, &slide)
		}
	}

	return slides, nil
}

func (db *slideRepo) InsertSlide(slide *model.Slide) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.conn.ExecContext(ctx, stmtInsertSlide, slide.Id, slide.PresentationId, slide.Type)
	if err != nil {
		return err
	}

	return nil
}

func (db *slideRepo) InsertContent(slideId string, content *model.Content) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.conn.ExecContext(ctx, stmtInsertContent, content.Id, slideId, content.Title, content.Meta)
	if err != nil {
		return err
	}

	return nil
}

func (db *slideRepo) InsertOption(contentId string, options []*model.Option) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	for _, option := range options {
		_, err := db.conn.ExecContext(ctx, stmtInsertOption, option.Name, option.Image, contentId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *slideRepo) UpdateSlide(presId string, slide model.Slide) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, stmtUpdateSlide, slide.Type, presId, slide.Id)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

func (db *slideRepo) UpdateContent(slideId string, content model.Content) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, stmtUpdateContent, content.Title, content.Meta, slideId)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

func (db *slideRepo) UpdateOptions(contentId string, options []*model.Option) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	for _, option := range options {
		_, err := db.conn.ExecContext(ctx, stmtUpdateOption, option.Name, option.Image, contentId, option.Id)
		if err != nil {
			return -1, err
		}
	}

	return 0, nil
}

func (db *slideRepo) UpdateOptionVote(contentId string, optionId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, stmtUpdateOptionVote, optionId, contentId)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

func (db *slideRepo) DeleteSlide(presId, slideId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, stmtDeleteSlideById, presId, slideId)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}
