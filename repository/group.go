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
	FindCreatedGroupsByUserId(userId string) ([]*model.Group, error)
	FindJoinedGroupsByUserId(userId string) ([]*model.GroupUser, error)
	FindGroupMemberDetailsByGroupId(groupId string) ([]*model.GroupUser, error)
	FindGroupById(groupId string) (*model.Group, error)
	UpdateUserRole(groupId, userId, role string) (int64, error)
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

	rows, err := db.conn.QueryContext(ctx, stmtSelectAllGroups)
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

	_, err = db.conn.ExecContext(ctx, stmtInsertGroupMember,
		id,
		userId,
		1,
		time.Now(),
	)

	if err != nil {
		return -1, err
	}

	group.Id = uint(id)
	return id, nil
}

func (db *groupRepo) FindCreatedGroupsByUserId(userId string) ([]*model.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx, stmtSelectGroupsByUserId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		var group model.Group
		if err = rows.Scan(
			&group.Id,
			&group.Name,
			&group.Link,
			&group.Desc,
			&group.CreatedAt); err != nil {
			return nil, errors.New("error scanning")
		}

		groups = append(groups, &group)
	}

	return groups, nil
}

func (db *groupRepo) FindJoinedGroupsByUserId(userId string) ([]*model.GroupUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx, stmtSelectJoinedGroupsByUserId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupUsers []*model.GroupUser
	for rows.Next() {
		var group model.Group
		var groupUser model.GroupUser

		if err = rows.Scan(
			&group.Id,
			&group.Name,
			&group.Link,
			&group.Desc,
			&group.CreatedAt,
			&groupUser.Role,
			&groupUser.JoinedAt); err != nil {
			return nil, errors.New("error scanning")
		}

		groupUser.GroupInfo = &group
		groupUsers = append(groupUsers, &groupUser)
	}

	return groupUsers, nil
}

func (db *groupRepo) FindGroupMemberDetailsByGroupId(groupId string) ([]*model.GroupUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx, stmtSelectGroupMemberDetailsById, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupUsers []*model.GroupUser
	for rows.Next() {
		var user model.User
		var groupUser model.GroupUser
		if err = rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Username,
			&user.SavedPassword,
			&user.Address,
			&user.ProfileImg,
			&user.Email,
			&user.CreatedAt,
			&groupUser.Role,
			&groupUser.JoinedAt); err != nil {
			return nil, errors.New("error scanning")
		}

		groupUser.UserInfo = &user
		groupUsers = append(groupUsers, &groupUser)
	}

	return groupUsers, nil
}

func (db *groupRepo) FindGroupById(groupId string) (*model.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var group model.Group
	var user model.User
	err := db.conn.QueryRowContext(ctx, stmtSelectGroupById, groupId).Scan(
		&group.Id,
		&group.Name,
		&group.Link,
		&group.Desc,
		&group.CreatedAt,
		&user.Id)

	if err != nil {
		return nil, errors.New("error scanning")
	}

	group.Owner = &user
	return &group, nil
}

func (db *groupRepo) UpdateUserRole(groupId, userId, role string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	updateResult, err := db.conn.ExecContext(ctx, stmtUpdateUserRole, role, groupId, userId)
	if err != nil {
		return -1, err
	}

	return updateResult.LastInsertId()
}
