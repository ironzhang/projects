CREATE DATABASE IF NOT EXISTS test;

CREATE TABLE IF NOT EXISTS db_test.tb_account(
	`type` int(8) NOT NULL,
	`name` varchar(64) NOT NULL,
);

CREATE TABLE IF NOT EXISTS db_test.tb_token(
	`id` bigint(20) NOT NULL AUTO_INCREMENT,
	`token` varchar(32) NOT NULL,
	`createAt` datetime NOT NULL,
	`modifyAt` datetime NOT NULL,
	PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
