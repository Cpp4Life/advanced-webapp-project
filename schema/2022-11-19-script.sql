-- INITIAL VERSION --

CREATE TABLE `users`
(
    `id`                BIGINT AUTO_INCREMENT PRIMARY KEY,
    `username`          VARCHAR(50)                NOT NULL,
    `password`          VARCHAR(250)               NOT NULL,
    `full_name`         VARCHAR(50)                NOT NULL,
    `email`             VARCHAR(50)                NOT NULL,
    `address`           VARCHAR(250) DEFAULT '',
    `profile_img`       VARCHAR(250) DEFAULT '',
    `user_tel`          VARCHAR(13)  DEFAULT '',
    `is_verified`       BOOLEAN      DEFAULT FALSE NOT NULL,
    `verification_code` VARCHAR(50)  DEFAULT '',
    `created_at`        DATETIME,
    `updated_at`        DATETIME
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE `groups`
(
    `id`         BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`       VARCHAR(50)  NOT NULL,
    `owner`      BIGINT       NOT NULL,
    `link`       VARCHAR(250) NOT NULL,
    `desc`       VARCHAR(250),
    `created_at` DATETIME
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE `roles`
(
    `id`    BIGINT AUTO_INCREMENT PRIMARY KEY,
    `title` VARCHAR(10) NOT NULL,
    `des`   VARCHAR(250)
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE `group_members`
(
    `member_id` BIGINT   NOT NULL,
    `group_id`  BIGINT   NOT NULL,
    `joined_at` DATETIME NOT NULL,
    `role`      BIGINT   NOT NULL,
    PRIMARY KEY (`member_id`, `group_id`)
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE `group_pres_infos`
(
    `group_id` BIGINT PRIMARY KEY,
    `pres_id`  BIGINT,
    `user_id`  BIGINT
) DEFAULT CHARSET = utf8mb4;

ALTER TABLE `groups`
    ADD(
        CONSTRAINT `groups_users_id_fk` FOREIGN KEY (`owner`) REFERENCES `users` (`id`)
        );

ALTER TABLE `group_members`
    ADD(
        CONSTRAINT `group_members_users_id_fk` FOREIGN KEY (`member_id`) REFERENCES `users` (`id`)
            ON DELETE CASCADE,
        CONSTRAINT `group_members_groups_id_fk` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`)
            ON DELETE CASCADE,
        CONSTRAINT `group_members_roles_id_fk` FOREIGN KEY (`role`) REFERENCES `roles` (`id`)
        );
