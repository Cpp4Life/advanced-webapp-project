package repository

import (
	"advanced-webapp-project/model"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IUserRepo interface {
	InsertUser(user *model.User) (int64, error)
	FindUserByEmail(email string) (*model.User, error)
	FindVerifiedStatusByEmail(email string) (bool, error)
	UpdateVerifiedStatus(verificationCode string) (int64, error)
	VerifyCredential(email, password string) (*model.User, error)
	FindUserById(id string) (*model.User, error)
	ModifyUserById(id string, user model.User) (int64, error)
}

type userRepo struct {
	conn *sql.DB
}

func NewUserRepo(sqldb *sql.DB) *userRepo {
	return &userRepo{
		conn: sqldb,
	}
}

func (db *userRepo) InsertUser(user *model.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return -1, err
	}

	user.SavedPassword = string(hashedPassword)

	insertResult, err := db.conn.ExecContext(ctx, stmtInsertUser,
		user.Username,
		user.SavedPassword,
		user.FullName,
		user.Address,
		user.ProfileImg,
		user.UserTel,
		user.Email,
		false,
		user.VerificationCode,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return -1, err
	}

	id, _ := insertResult.LastInsertId()
	user.Id = uint(id)
	user.Password = ""
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
			&user.SavedPassword,
			&user.Address,
			&user.ProfileImg,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *userRepo) FindVerifiedStatusByEmail(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var isVerified bool
	err := db.conn.QueryRowContext(ctx, stmtSelectVerifiedStatusByEmail, email).Scan(&isVerified)
	if err != nil {
		return false, err
	}

	return isVerified, nil
}

func (db *userRepo) UpdateVerifiedStatus(verificationCode string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := db.conn.ExecContext(ctx, stmtUpdateVerifiedStatus, verificationCode)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (db *userRepo) VerifyCredential(email, password string) (*model.User, error) {
	user, err := db.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.SavedPassword), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (db *userRepo) FindUserById(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user model.User
	err := db.conn.QueryRowContext(ctx, stmtSelectUserById, id).
		Scan(
			&user.Id,
			&user.FullName,
			&user.Username,
			&user.SavedPassword,
			&user.Address,
			&user.ProfileImg,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *userRepo) ModifyUserById(id string, user model.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	updateResult, err := db.conn.ExecContext(ctx, stmtUpdateUserById,
		user.FullName,
		user.Username,
		user.ProfileImg,
		time.Now(),
		id,
	)
	if err != nil {
		return -1, err
	}

	return updateResult.RowsAffected()
}
