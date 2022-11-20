package repository

import "time"

const (
	dbTimeout = 5 * time.Second

	stmtInsertUser = "INSERT INTO `users` " +
		"(username, password, full_name, address, profile_img, user_tel, email, created_at, updated_at) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);"

	stmtSelectUserByEmail = "SELECT email FROM `users` " +
		"WHERE email LIKE ?;"
)
