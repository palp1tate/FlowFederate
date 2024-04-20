CREATE TABLE `user_info` (
  `uuid` varchar(255) NOT NULL,
  `user_name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` int NOT NULL DEFAULT '1' COMMENT '0为管理员，1为普通用户',
  `state` int NOT NULL DEFAULT '0' COMMENT '0为启用，1为禁用',
  `create_time` varchar(255) NOT NULL,
  PRIMARY KEY (`user_name`) USING BTREE,
  CONSTRAINT `user_info_chk_1` CHECK ((`role` in (0,1))),
  CONSTRAINT `user_info_chk_2` CHECK ((`state` in (0,1)))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `user_info` (`uuid`, `user_name`, `password`, `role`, `state`, `create_time`) VALUES ('db9d423c791f11ee93240242ac110004', 'admin', '123456', 0, 0, '2023-11-02 01:33:43');
INSERT INTO `user_info` (`uuid`, `user_name`, `password`, `role`, `state`, `create_time`) VALUES ('613bc48f7f7311ee8a130242ac110005', 'test', '123456', 1, 0, '2023-11-10 10:46:42');
