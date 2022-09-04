CREATE TABLE IF NOT EXISTS `users`
(
    `id`      bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` int(11) NOT NULL,
    `name`    varchar(64) DEFAULT '' NOT NULL,
    PRIMARY KEY (`id`)
);
-- user1 user2
CREATE TABLE IF NOT EXISTS `friend_link`
(
    `id`       bigint(20) NOT NULL AUTO_INCREMENT,
    `user1_id` int(11) NOT NULL,
    `user2_id` int(11) NOT NULL,
    PRIMARY KEY (`id`)
);
-- user1 user2 block
CREATE TABLE IF NOT EXISTS `block_list`
(
    `id`       bigint(20) NOT NULL AUTO_INCREMENT,
    `user1_id` int(11) NOT NULL,
    `user2_id` int(11) NOT NULL,
    PRIMARY KEY (`id`)
);
