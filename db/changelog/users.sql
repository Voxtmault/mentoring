--liquibase formatted sql

--changeset Voxtmault:1
CREATE TABLE IF NOT EXISTS `users` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `password` LONGTEXT NOT NULL,
    `salt` LONGTEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME ON UPDATE CURRENT_TIMESTAMP NULL DEFAULT NULL,
    `deleted_at` DATETIME NULL DEFAULT NULL
)ENGINE = InnoDB;
--rollback DROP TABLE `users`;