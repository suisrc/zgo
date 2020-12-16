-- 用户
INSERT INTO `zgo_user`(`id`, `kid`, `name`) VALUES
(1, "100001", "zgo-user"),
(2, "100002", "plus-user"),
(3, "100003", "suisrc-user"),
(4, "100004", "admin-user");

-- 账户
INSERT INTO `zgo_account`(`id`, `account`, `account_typ`, `password`, `password_salt`, `password_type`, `user_id`, `role_id`) VALUES 
(1, "zgo1",        1, "c557193f596ccf70b8cbc5ca306557b3", "uoqacs2t699ybv8tc42hz8z1shny6ups", "MD5", 1, null),
(2, "zgo2",        1, "654321", null, null, 1, null),
(3, "plus1",       1, "123456", null, null, 2, null),
(4, "ss",          1, "uBnKfXylWRdUFqVM424ERH.tISbfJbWq", "J3Apb1ZhNgtuBx4ifhg9F0MBVhI3bH9ELjJRQg==", "BCR3", 3, null),
(5, "admin",       1, "123456", null, null, 4, null),
(6, "13311111111", 2, null,     null, null, 4, null);

-- domain=jwt代表特殊含义, 系统需要访问的domain和jwt中的aud是否相等
INSERT INTO `zgo_gateway`(`id`, `name`, `domain`, `methods`, `path`, `netmask`, `allow`, `status`) VALUES 
(1, "nosignin", null,                               null,                           "/*",        null,          0, 1),
(2, "norole",   null,                               null,                           "/*",        null,          0, 1),
(3, "admin",    null,                               null,                           "/*",        null,          1, 1),
(4, "api-demo", "xxx.api.io",                       "(GET)|(POST)|(PUT)|(DELETE)",  "/api/*",    "0.0.0.0/0",   1, 1),
(5, "jwt",      "jwt",                              null,                           "/*",        null,          1, 1),
(6, "lys2go",   "leo2go-80-me03.dev.sims-cn.com",   null,                           "/*",        null,          1, 1);

-- domain 如果同一用户有多个角色,系统可以通过访问的域名自动分配角色
INSERT INTO `zgo_role`(`id`, `kid`, `name`, `domain`, `status`, `organization`) VALUES 
(1, "admin",    "管理员",  null,                             1, null),
(2, "normal",   "正常",    null,                             1, null),
(3, "group1",   "分组一",  null,                             1, null),
(4, "group2",   "分组二",  null,                             1, null),
(5, "group3",   "分组三",  "leo2go-80-me03.dev.sims-cn.com", 1, null),
(6, "invalid",  "作废",    null,                             0, null);

INSERT INTO `zgo_role_gateway`(`id`, `role_id`, `gateway`) VALUES 
(1, 1, "admin"),
(2, 2, "api"),
(3, 2, "jwt"),
(4, 6, "lys2go");

INSERT INTO `zgo_role_role`(`id`, `owner_id`, `child_id`) VALUES 
(1, 3, 1),
(2, 4, 2),
(3, 5, 6);

INSERT INTO `zgo_user_role`(`id`, `user_id`, `role_id`) VALUES 
(1, 1, 1),
(2, 2, 2),
(3, 3, 3),
(4, 3, 4),
(5, 4, 5),
(6, 4, 6);

INSERT INTO `zgo_oauth2_client` (`id`, `kid`, `audience`, `issuer`, `expired`, `token_secret`) VALUES
(1, "crm", "crm.dev.api.io", ".dev.icgear.cn", 3600, "11111111"),
(2, "sto", "sto.dev.api.io", ".dev.icgear.cn", 1800, "22222222");
