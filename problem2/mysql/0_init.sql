DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`      bigint(20) unsigned    NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned       NOT NULL UNIQUE,
    `name`    varchar(64) DEFAULT '' NOT NULL,
    PRIMARY KEY (`id`)
);
-- user1 user2
DROP TABLE IF EXISTS `friend_link`;
CREATE TABLE `friend_link`
(
    `id`       bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user1_id` int(11) unsigned    NOT NULL,
    `user2_id` int(11) unsigned    NOT NULL,
    PRIMARY KEY (`id`)
);
-- user1 user2 block
DROP TABLE IF EXISTS `block_list`;
CREATE TABLE `block_list`
(
    `id`       bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user1_id` int(11) unsigned    NOT NULL,
    `user2_id` int(11) unsigned    NOT NULL,
    PRIMARY KEY (`id`)
);
