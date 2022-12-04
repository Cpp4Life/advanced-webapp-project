INSERT INTO `roles` VALUES (1, 'owner', 'user has full privileged permission to operate in group');
INSERT INTO `roles` VALUES (2, 'co-owner', 'user has the same permission as `owner` user');
INSERT INTO `roles` VALUES (3, 'member', 'user has limited permissions in group');

INSERT INTO `question_categories`(`name`) VALUES ('Popular question types');
INSERT INTO `question_categories`(`name`) VALUES ('Quiz competition');
INSERT INTO `question_categories`(`name`) VALUES ('Content slides');
INSERT INTO `question_categories`(`name`) VALUES ('Advanced questions');

INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Multiple Choice', 1);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Word Cloud', 1);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Open Ended', 1);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Scales', 1);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Ranking', 1);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Q&A', 1);

INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Select Answer', 2);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Type Answer', 2);

INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Heading', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Paragraph', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Bullets', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Image', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Video', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Big', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Quote', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Number', 3);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Instructions', 3);

INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('100 points', 4);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('2x2 Grid', 4);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Who will win?', 4);
INSERT INTO `question_types`(`name`, `question_cate_id`) VALUES('Pin on Image', 4);