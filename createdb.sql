CREATE TABLE `user` (
    `id` integer,
    `username` text UNIQUE,
    `hash` text,
    PRIMARY KEY(`id`)
);

CREATE TABLE `alias` (
    `id` integer,
	`name`	text UNIQUE,
	`url`	text,
    `userid` integer,
	PRIMARY KEY(`id`)
    FOREIGN KEY(`userid`) REFERENCES `user`(`id`)
);

INSERT INTO `user` VALUES (0, 'nouser', '');
