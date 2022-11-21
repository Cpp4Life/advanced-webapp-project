package repository

import (
	"advanced-webapp-project/model"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IUserRepo interface {
	InsertUser(user model.User) (int64, error)
	FindUserByEmail(email string) (*model.User, error)
	VerifyCredential(email, password string) (*model.User, error)
}

type userRepo struct {
	conn *sql.DB
}

func NewUserRepo(sqldb *sql.DB) *userRepo {
	return &userRepo{
		conn: sqldb,
	}
}

func (db *userRepo) InsertUser(user model.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return -1, err
	}

	insertResult, err := db.conn.ExecContext(ctx, stmtInsertUser,
		user.Username,
		hashedPassword,
		user.FullName,
		user.Address,
		user.ProfileImg,
		user.UserTel,
		user.Email,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return -1, err
	}

	id, _ := insertResult.LastInsertId()
	return id, nil
}

func (db *userRepo) FindUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user model.User
	err := db.conn.QueryRowContext(ctx, stmtSelectUserByEmail, email).
		Scan(
			&user.Id,
			&user.FullName,
			&user.Username,
			&user.Password,
			&user.FullName,
			&user.Address,
			&user.ProfileImg,
			&user.Email,
		)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *userRepo) VerifyCredential(email, password string) (*model.User, error) {
	user, err := db.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
