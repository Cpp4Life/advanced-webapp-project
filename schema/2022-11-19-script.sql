-- INITIAL VERSION --

CREATE TABLE `users` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(50) NOT NULL,
    `password` VARCHAR(50) NOT NULL,
    `full_name` VARCHAR(50) NOT NULL,
    `address` VARCHAR(250),
    `profile_img` VARCHAR(250),
    `user_tel` VARCHAR(13)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `groups` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(50) NOT NULL,
    `link` VARCHAR(100) NOT NULL,
    `desc` VARCHAR(250),
    `created_date` DATETIME NOT NULL,
    `owner` BIGINT REFERENCES users(`id`)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `roles` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `title` VARCHAR(10) NOT NULL,
    `des` VARCHAR(250)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `group_members` (
    `member_id` BIGINT REFERENCES users(`id`),
    `group_id` BIGINT REFERENCES groups(`id`),
    `joined_at` DATETIME NOT NULL,
    `role` VARCHAR(10) REFERENCES roles(`id`),
    PRIMARY KEY (`member_id`, `group_id`)
) DEFAULT CHARSET=utf8mb4;