CREATE TABLE `user` (
                        `id` int(11) NOT NULL AUTO_INCREMENT,
                        `add_time` datetime(3) DEFAULT NULL,
                        `update_time` datetime(3) DEFAULT NULL,
                        `delete_time` datetime(3) DEFAULT NULL,
                        `is_deleted` tinyint(1) DEFAULT NULL,
                        `email` varchar(100) DEFAULT NULL,
                        `password` varchar(100) DEFAULT NULL,
                        `nick_name` varchar(100) DEFAULT NULL,
                        `birthday` datetime DEFAULT NULL,
                        `gender` varchar(6) DEFAULT 'male',
                        `role` bigint(20) DEFAULT '1',
                        `desc` text,
                        `image` varchar(200) DEFAULT NULL,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `uni_user_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4


