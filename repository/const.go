package repository

import "time"

const (
	dbTimeout = 5 * time.Second

	stmtInsertUser = "INSERT INTO `users` " +
		"(username, password, full_name, address, profile_img, user_tel, email, is_verified, verification_code, created_at, updated_at) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"

	stmtSelectUserByEmail = "SELECT id, full_name, username, password, address, profile_img, email, created_at, updated_at " +
		"FROM `users` " +
		"WHERE email LIKE ?;"

	stmtSelectUserById = "SELECT id, full_name, username, password, address, profile_img, email, created_at, updated_at " +
		"FROM `users` " +
		"WHERE id = ?;"

	stmtSelectVerifiedStatusByEmail = "SELECT is_verified " +
		"FROM `users` " +
		"WHERE email LIKE ?;"

	stmtUpdateVerifiedStatus = "UPDATE `users` " +
		"SET is_verified = true " +
		"WHERE verification_code LIKE ?;"

	stmtUpdateUserById = "UPDATE `users` " +
		"SET full_name = ?, username = ?, profile_img = ?, updated_at = ? " +
		"WHERE id = ?;"

	stmtInsertGroup = "INSERT INTO `groups` " +
		"(name, link, `desc`, created_at, owner) " +
		"VALUES (?, ?, ?, ?, ?);"

	stmtInsertGroupMember = "INSERT INTO `group_members` " +
		"(group_id, member_id, role, joined_at) " +
		"VALUES (?, ?, ?, ?);"

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

	stmtUpdateUserRole = "UPDATE `group_members` " +
		"SET role = ? " +
		"WHERE group_id = ? AND member_id = ? "
)
