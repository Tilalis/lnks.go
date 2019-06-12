CREATE TABLE `user` (
    `id` integer,
    `username` text UNIQUE,
    `hash` text,
    PRIMARY KEY(`id`)
);

CREATE TABLE `alias` (
	`name`	text,
	`url`	text,
    `userid` integer,
	PRIMARY KEY(`name`)
    FOREIGN KEY(`userid`) REFERENCES `user`(`id`)
);

INSERT INTO `user` VALUES (0, 'nouser', '');
