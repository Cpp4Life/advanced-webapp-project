package db

import (
	"advanced-webapp-project/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewSQLDB() *sql.DB {
	return connect()
}

func getSourceString() string {
	appConfig, _ := config.LoadConfig(".")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		appConfig.DBUser,
		appConfig.DBPass,
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBName,
	)

	return dsn
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", getSourceString())
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)

	return db
}

func Close(sqlDB *sql.DB) {
	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}
