CREATE SCHEMA IF NOT EXISTS `registration` ;
USE `registration` ;
CREATE TABLE `user` (
    `id` varchar(255),
    `firstname` varchar(255),
    `lastname` varchar(255),
    `email` varchar(255),
    `password` varchar(255)
);

CREATE TABLE `blacklist` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `token` varchar(255)
);