package repository

import (
	"advanced-webapp-project/model"
	"context"
	"database/sql"
	"errors"
	"time"
)

type IGroupRepo interface {
	FindAll() ([]*model.Group, error)
	InsertGroup(group *model.Group, userId string) (int64, error)
}

type groupRepo struct {
	conn *sql.DB
}

func NewGroupRepo(sqldb *sql.DB) *groupRepo {
	return &groupRepo{
		conn: sqldb,
	}
}

func (db *groupRepo) FindAll() ([]*model.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx, stmtFindAllGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		var group model.Group
		var user model.User
		if err = rows.Scan(
			&group.Id,
			&group.Name,
			&group.Link,
			&group.Desc,
			&group.CreatedAt,
			&user.Id); err != nil {
			return nil, errors.New("error scanning")
		}
		group.Owner = &user
		groups = append(groups, &group)
	}

	return groups, nil
}

func (db *groupRepo) InsertGroup(group *model.Group, userId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	insertResult, err := db.conn.ExecContext(ctx, stmtInsertGroup,
		group.Name,
		group.Link,
		group.Desc,
		time.Now(),
		userId,
	)

	if err != nil {
		return -1, err
	}

	id, _ := insertResult.LastInsertId()
	group.Id = uint(id)
	return id, nil
}
