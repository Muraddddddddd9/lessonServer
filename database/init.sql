-- MySQL database schema for Lesson Server
-- Version: 1.0
-- Created: 2025-07-05

CREATE DATABASE IF NOT EXISTS `lesson_server` 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE `lesson_server`;

DROP TABLE IF EXISTS `setting`;
CREATE TABLE `setting` (
  `id` int NOT NULL AUTO_INCREMENT,
  `now_stage_lesson` varchar(45) NOT NULL DEFAULT '0',
  `id_presentation` varchar(45) NOT NULL DEFAULT '',
  `test_team_first` varchar(45) NOT NULL DEFAULT '',
  `test_team_second` varchar(45) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `password` varchar(255) NOT NULL,
  `status` varchar(45) NOT NULL,
  `bim_coin1` int NOT NULL DEFAULT '0',
  `bim_coin2` float NOT NULL DEFAULT '0',
  `team` int NOT NULL,
  `test_first` tinyint(1) NOT NULL DEFAULT '0',
  `time_test` varchar(45) NOT NULL DEFAULT '0',
  `bim_total` float NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `answer_test_two`;
CREATE TABLE `answer_test_two` (
  `team` int NOT NULL,
  `conflict_elements` text,
  `visualization1` text,
  `visualization2` text,
  `visualization3` text,
  `matrix_file` text,
  `collisions_count` text,
  `priority` text,
  `risk_assessment` text,
  `consequences` text,
  `problems` text,
  `specialists` text,
  `stage` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`team`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;