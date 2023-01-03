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

	stmtUpdatePassword = "UPDATE `users` " +
		"SET password = ?, updated_at = ? " +
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

	stmtSelectGroupPresentationInfo = "SELECT group_id, pres_id, user_id " +
		"FROM `group_pres_infos` " +
		"WHERE group_id = ?;"

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

	stmtDeleteGroupById = "DELETE FROM `groups` WHERE id = ?;"

	stmtDeleteAllGroupMembers = "DELETE FROM `group_members` WHERE group_id = ?;"

	stmtSelectUserRole = "SELECT `role` " +
		"FROM `group_members` " +
		"WHERE group_id = ? AND member_id = ?;"

	stmtDeleteMember = "DELETE FROM `group_members` WHERE group_id = ? AND member_id = ?;"

	stmtSelectPresentationById = "SELECT id, name, owner, modified_at, created_at " +
		"FROM `presentations` WHERE id = ?; "

	stmtSelectAllPresentationsByUserId = "SELECT id, name, modified_at, created_at " +
		"FROM `presentations` " +
		"WHERE owner = ?;"

	stmtInsertPresentation = "INSERT INTO `presentations` " +
		"(id, name, owner, modified_at, created_at) " +
		"VALUES (?, ?, ?, ?, ?);"

	stmtUpdatePresentation = "UPDATE `presentations` " +
		"SET name = ?, modified_at = ? " +
		"WHERE id = ?; "

	stmtDeletePresentation = "DELETE `p`, `s`, `c` " +
		"FROM `presentations` p " +
		"JOIN slides s on p.id = s.pres_id " +
		"JOIN contents c on s.id = c.slide_id " +
		"WHERE p.id = ?;"

	stmtReplaceGroupPresentation = "REPLACE INTO `group_pres_infos` VALUES(?, ?, ?)"

	stmtInsertSlide = "INSERT INTO `slides` " +
		"(id, pres_id, slide_type) " +
		"VALUES (?, ?, ?);"

	stmtInsertContent = "INSERT INTO `contents` " +
		"(id, slide_id, title, meta) " +
		"VALUES (?, ?, ?, ?);"

	stmtInsertOption = "INSERT INTO `options` " +
		"(name, image, content_id) " +
		"VALUES (?, ?, ?);"

	stmtInsertHeading = "INSERT INTO `headings` " +
		"(heading, sub_heading, image, content_id) " +
		"VALUES (?, ?, ?, ?);"

	stmtInsertParagraph = "INSERT INTO `paragraphs` " +
		"(heading, text, image, content_id) " +
		"VALUES (?, ?, ?, ?);"

	stmtUpdateSlide = "UPDATE `slides` " +
		"SET slide_type = ? " +
		"WHERE pres_id = ? AND id = ?;"

	stmtUpdateContent = "UPDATE `contents` " +
		"SET title = ?, meta = ? " +
		"WHERE slide_id = ?;"

	stmtUpdateOption = "UPDATE `options` " +
		"SET name = ?, image = ? " +
		"WHERE content_id = ? AND id = ?;"

	stmtUpdateHeading = "UPDATE `headings` " +
		"SET heading = ?, sub_heading = ?, image = ? " +
		"WHERE content_id = ? AND id = ?;"

	stmtUpdateParagraph = "UPDATE `paragraphs` " +
		"SET heading = ?, text = ?, image = ? " +
		"WHERE content_id = ? AND id = ?;"

	stmtSelectAllSlides = "WITH `sub-contents`(id, heading, sub_heading, image, total_votes, content_id) AS " +
		"( " +
		"    SELECT o.id, o.name, '', o.image, o.total_votes, content_id " +
		"    FROM `options` o " +
		"    UNION " +
		"    SELECT h.id, h.heading, h.sub_heading, h.image, '', content_id " +
		"    FROM `headings` h " +
		"    UNION " +
		"    SELECT p.id, p.heading, p.text, p.image, '', content_id " +
		"    FROM `paragraphs` p " +
		") " +
		"SELECT s.id, s.slide_type, c.id, c.title, c.meta, sc.id, sc.heading, sc.sub_heading, sc.image, sc.total_votes " +
		"FROM slides s " +
		"JOIN contents c on s.id = c.slide_id " +
		"JOIN `sub-contents` sc on c.id = sc.content_id " +
		"WHERE s.pres_id = ?;"

	stmtSelectSlideById = "WITH `sub-contents`(id, heading, sub_heading, image, total_votes, content_id) AS " +
		"( " +
		"    SELECT o.id, o.name, '', o.image, o.total_votes, content_id " +
		"    FROM `options` o " +
		"    UNION " +
		"    SELECT h.id, h.heading, h.sub_heading, h.image, '', content_id " +
		"    FROM `headings` h " +
		"    UNION " +
		"    SELECT p.id, p.heading, p.text, p.image, '', content_id " +
		"    FROM `paragraphs` p " +
		") " +
		"SELECT s.id, s.slide_type, c.id, c.title, c.meta, sc.id, sc.heading, sc.sub_heading, sc.image, sc.total_votes " +
		"FROM slides s " +
		"JOIN contents c on s.id = c.slide_id " +
		"JOIN `sub-contents` sc on c.id = sc.content_id " +
		"WHERE s.id = ?;"

	stmtUpdateOptionVote = "UPDATE `options` " +
		"SET `total_votes` = `total_votes` + 1 " +
		"WHERE id = ? AND content_id = ?;"

	stmtDeleteSlideById = "DELETE FROM `slides` WHERE pres_id = ? AND id = ?;"
)
