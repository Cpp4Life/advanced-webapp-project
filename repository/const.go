package repository

import "time"

const (
	dbTimeout = 10 * time.Second

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

	stmtSelectPresentationById = "SELECT id, name, owner, modified_at, created_at " +
		"FROM `presentations` WHERE id = ?; "

	stmtSelectAllPresentations = "SELECT id, name, owner, modified_at, created_at FROM `presentations`; "

	stmtInsertPresentation = "INSERT INTO `presentations` " +
		"(name, owner, modified_at, created_at) " +
		"VALUES (?, ?, ?, ?);"

	stmtUpdatePresentation = "UPDATE `presentations` " +
		"SET name = ?, modified_at = ? " +
		"WHERE id = ?; "

	stmtDeletePresentation = "DELETE FROM `presentations` WHERE id = ?;"

	stmtInsertSlide = "INSERT INTO `slides` " +
		"(pres_id, slide_type) " +
		"VALUES (?, ?);"

	stmtInsertContent = "INSERT INTO `contents` " +
		"(slide_id, title, meta) " +
		"VALUES (?, ?, ?);"

	stmtInsertOption = "INSERT INTO `options` " +
		"(name, image, content_id) " +
		"VALUES (?, ?, ?);"

	stmtUpdateSlide = "UPDATE `slides` " +
		"SET slide_type = ? " +
		"WHERE pres_id = ? AND id = ?;"

	stmtUpdateContent = "UPDATE `contents` " +
		"SET title = ?, meta = ? " +
		"WHERE slide_id = ?;"

	stmtUpdateOption = "UPDATE `options` " +
		"SET name = ?, image = ? " +
		"WHERE content_id = ? AND id = ?;"

	stmtSelectAllSlides = "SELECT s.id, s.slide_type, c.id, c.title, c.meta, o.id, o.name, o.image " +
		"FROM `slides` s " +
		"JOIN `contents` c on s.id = c.slide_id " +
		"JOIN `options` o on c.id = o.content_id " +
		"WHERE s.pres_id = ?;"
)
