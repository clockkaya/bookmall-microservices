-- `go`.`user` definition

CREATE TABLE `user` (
                        `id` bigint NOT NULL AUTO_INCREMENT,
                        `number` varchar(255) NOT NULL DEFAULT '' COMMENT '学号',
                        `name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名称',
                        `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码',
                        `gender` char(5) NOT NULL COMMENT '男｜女｜未公开',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        `points` bigint NOT NULL DEFAULT '100' COMMENT '用户积分',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `number_unique` (`number`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;