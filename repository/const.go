package repository

import "time"

const (
	dbTimeout = 5 * time.Second

	stmtInsertUser = "INSERT INTO `users` " +
		"(username, password, full_name, address, profile_img, user_tel, email, created_at, updated_at) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);"

	stmtSelectUserByEmail = "SELECT id, full_name, username, password, address, profile_img, email, created_at " +
		"FROM `users` " +
		"WHERE email LIKE ?;"

	stmtSelectUserById = "SELECT id, full_name, username, password, address, profile_img, email, created_at " +
		"FROM `users` " +
		"WHERE id = ?;"

	stmtInsertGroup = "INSERT INTO `groups` " +
		"(name, link, `desc`, created_at, owner) " +
		"VALUES (?, ?, ?, ?, ?);"

	stmtFindAllGroups = "SELECT id, name, link, `desc`, created_at, owner FROM `groups`;"
)
