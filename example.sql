CREATE TABLE `t_demo`(
    `id`                                int(9) NOT NULL AUTO_INCREMENT,
    `uid`                               mediumint(9) NOT NULL COMMENT 'uid',
    `create_time`                       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT='demo测试表';