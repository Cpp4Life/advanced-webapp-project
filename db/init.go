package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type sqlDB struct {
	sqlDB *sql.DB
}

func NewSQLDB() *sqlDB {
	return &sqlDB{
		sqlDB: connect(),
	}
}

func getSourceString() string {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return dsn
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", getSourceString())
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func (i *sqlDB) Close() {
	if err := i.sqlDB.Close(); err != nil {
		panic(err)
	}
}
