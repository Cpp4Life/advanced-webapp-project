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

	stmtUpdateUserById = "UPDATE `users` " +
		"SET full_name = ?, username = ?, profile_img = ? " +
		"WHERE id = ?;"

	stmtInsertGroup = "INSERT INTO `groups` " +
		"(name, link, `desc`, created_at, owner) " +
		"VALUES (?, ?, ?, ?, ?);"

	stmtSelectAllGroups = "SELECT id, name, link, `desc`, created_at, owner FROM `groups`;"

	stmtSelectGroupById = "SELECT id, name, link, `desc`, created_at, owner " +
		"FROM `groups` " +
		"WHERE id = ?;"

	stmtSelectGroupsByUserId = "SELECT id, name, link, `desc`, created_at " +
		"FROM `groups` " +
		"WHERE owner = ?;"

	stmtSelectJoinedGroupsByUserId = "SELECT g.id, g.name, g.link, g.`desc`, g.created_at, r.title, joined_at " +
		"FROM `group_members` " +
		"JOIN `groups` g ON group_members.group_id = g.id " +
		"JOIN `roles` r ON group_members.role = r.id " +
		"WHERE `member_id` = ?;"

	stmtSelectGroupMemberDetailsById = "SELECT u.id, u.full_name, u.username, u.password, u.address, u.profile_img, u.email, u.created_at, r.title, joined_at " +
		"FROM `group_members` " +
		"JOIN `users` u ON group_members.member_id = u.id " +
		"JOIN `roles` r ON group_members.role = r.id " +
		"WHERE `group_id` = ?;"
)
