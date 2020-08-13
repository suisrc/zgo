INSERT INTO `user`(`id`, `uid`, `name`) VALUES
(1, "100001", "zgo"),
(2, "100002", "icgear"),
(3, "100003", "suisrc"),
(4, "100004", "user");

INSERT INTO `account`(`id`, `account`, `password`, `password_salt`, `password_type`, `user_id`, `role_id`, `oauth2_id`) VALUES 
(1, "zgo",  "c557193f596ccf70b8cbc5ca306557b3", "uoqacs2t699ybv8tc42hz8z1shny6ups", "MD5", 1, null, null),
(2, "zgo2", "654321", null, null, 1, null, null),
(3, "icg",  "123456", null, null, 2, null, null),
(4, "ss",   "uBnKfXylWRdUFqVM424ERH.tISbfJbWq", "J3Apb1ZhNgtuBx4ifhg9F0MBVhI3bH9ELjJRQg==", "BCR3", 3, null, null),
(5, "user", "123456", null, null, 4, null, null);

INSERT INTO `resource`(`id`, `resource`, `domain`, `methods`, `path`, `netmask`, `allow`, `status`) VALUES 
(1, "nosignin", null,         null,                          "/*",        null,       0, 1),
(2, "norole",   null,         null,                          "/*",        null,       1, 1),
(3, "admin",    null,         null,                          "/*",        null,       1, 1),
(4, "api",      "xxx.api.io", "(GET)|(POST)|(PUT)|(DELETE)", "/api/xxx", "0.0.0.0/0", 1, 1),
(5, "jwt",      "jwt",        null,                          "/*",       null,        1, 1);

INSERT INTO `role`(`id`, `uid`, `name`, `status`) VALUES 
(1, "admin",    "管理员", 1),
(2, "normal",   "正常",   1),
(3, "group1",   "分组一", 1),
(4, "group2",   "分组二", 1),
(5, "group3",   "分组三", 1),
(6, "invalid",  "作废",   0);

INSERT INTO `resource_role`(`id`, `role_id`, `resource`) VALUES 
(1, 1, "admin"),
(2, 2, "api"),
(3, 2, "jwt"),
(4, 6, "admin");

INSERT INTO `role_role`(`id`, `owner_id`, `child_id`) VALUES 
(1, 3, 1),
(2, 4, 2),
(3, 5, 6);

INSERT INTO `user_role`(`id`, `user_id`, `role_id`) VALUES 
(1, 1, 1),
(2, 2, 2),
(3, 3, 3),
(4, 3, 4),
(5, 4, 5),
(6, 4, 6);