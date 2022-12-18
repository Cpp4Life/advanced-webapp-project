package repository

import (
	"advanced-webapp-project/model"
	"advanced-webapp-project/utils"
	"context"
	"database/sql"
)

type ISlideRepo interface {
	FindAllSlides(presId string) ([]*model.Slide, error)
	InsertSlide(slide *model.Slide) error
	InsertContent(slideId string, content *model.Content) error
	InsertOption(contentId string, options []*model.Option) error
	InsertHeading(contentId string, heading *model.Heading) error
	InsertParagraph(contentId string, paragraph *model.Paragraph) error
	UpdateSlide(presId string, slide model.Slide) (int64, error)
	UpdateContent(slideId string, content model.Content) (int64, error)
	UpdateOptions(contentId string, options []*model.Option) (int64, error)
	UpdateOptionVote(contentId string, optionId string) (int64, error)
	UpdateHeading(contentId string, heading *model.Heading) (int64, error)
	UpdateParagraph(contentId string, paragraph *model.Paragraph) (int64, error)
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

	type subContent struct {
		Id         uint
		Heading    string
		SubHeading string
		Image      string
		TotalVotes string
	}

	var slides []*model.Slide
	for rows.Next() {
		var slide model.Slide
		var content model.Content
		var option model.Option
		var heading model.Heading
		var paragraph model.Paragraph
		var sc subContent
		var index = len(slides) - 1
		if err = rows.Scan(
			&slide.Id,
			&slide.Type,
			&content.Id,
			&content.Title,
			&content.Meta,
			&sc.Id,
			&sc.Heading,
			&sc.SubHeading,
			&sc.Image,
			&sc.TotalVotes,
		); err != nil {
			return nil, err
		}

		switch slide.Type {
		case 1:
			option.Id = sc.Id
			option.Name = sc.Heading
			option.Image = sc.Image
			option.TotalVotes = utils.Str2Uint(sc.TotalVotes)
			content.Options = append(content.Options, &option)
		case 8:
			heading.Id = sc.Id
			heading.Heading = sc.Heading
			heading.SubHeading = sc.SubHeading
			heading.Image = sc.Image
			content.Heading = &heading
		case 9:
			paragraph.Id = sc.Id
			paragraph.Heading = sc.Heading
			paragraph.Text = sc.SubHeading
			paragraph.Image = sc.Image
			content.Paragraph = &paragraph
		default:
		}

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

func (db *slideRepo) InsertHeading(contentId string, heading *model.Heading) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.conn.ExecContext(ctx, stmtInsertHeading, heading.Heading, heading.SubHeading, heading.Image, contentId)
	if err != nil {
		return err
	}

	return nil
}

func (db *slideRepo) InsertParagraph(contentId string, paragraph *model.Paragraph) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.conn.ExecContext(ctx, stmtInsertParagraph, paragraph.Heading, paragraph.Text, paragraph.Image, contentId)
	if err != nil {
		return err
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

func (db *slideRepo) UpdateHeading(contentId string, heading *model.Heading) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, stmtUpdateHeading, heading.Heading, heading.SubHeading, heading.Image, contentId, heading.Id)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

func (db *slideRepo) UpdateParagraph(contentId string, paragraph *model.Paragraph) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, stmtUpdateParagraph, paragraph.Heading, paragraph.Text, paragraph.Image, contentId, paragraph.Id)
	if err != nil {
		return -1, nil
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
