-- MySQL database schema for Lesson Server
-- Version: 1.0
-- Created: 2025-07-05

CREATE DATABASE IF NOT EXISTS `lesson_server` 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE `lesson_server`;

DROP TABLE IF EXISTS `setting`;
CREATE TABLE `setting` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `now_stage_lesson` varchar(45) NOT NULL DEFAULT '0',
  `id_presentation` varchar(45) NOT NULL DEFAULT '',
  `test_team_first` tinyint(1) NOT NULL DEFAULT 0,
  `test_team_second` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `password` varchar(255) NOT NULL,
  `status` varchar(45) NOT NULL,
  `bim_coin` int NOT NULL DEFAULT '0',
  `team` int NOT NULL,
  `test_first` tinyint(1) NOT NULL DEFAULT '0',
  `time_test` varchar(45) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci