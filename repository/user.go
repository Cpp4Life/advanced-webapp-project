package repository

import (
	"advanced-webapp-project/model"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IUserRepo interface {
	InsertUser(user model.User) error
	FindUserByEmail(email string) (*model.User, error)
}

type userRepo struct {
	conn *sql.DB
}

func NewUserRepo(sqldb *sql.DB) *userRepo {
	return &userRepo{
		conn: sqldb,
	}
}

func (db *userRepo) InsertUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	_, err = db.conn.ExecContext(ctx, stmtInsertUser,
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
		return err
	}

	return nil
}

func (db *userRepo) FindUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user model.User
	err := db.conn.QueryRowContext(ctx, stmtSelectUserByEmail, email).Scan(&user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
