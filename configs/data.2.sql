-- 用户
INSERT INTO `zgo_user`(`id`, `kid`, `name`) VALUES
(1, "100001", "it-0001");

-- 账户
INSERT INTO `zgo_account`(`id`, `account`, `account_typ`, `password`, `password_salt`, `password_type`, `user_id`, `role_id`) VALUES 
(1, "it-0001", 1, "123456", null, null, 1, null);

-- domain=jwt代表特殊含义, 系统需要访问的domain和jwt中的aud是否相等
INSERT INTO `zgo_gateway`(`id`, `name`, `domain`, `methods`, `path`, `netmask`, `allow`, `status`) VALUES 
(1, "nosignin", null,  null, "/*", null, 0, 1),
(2, "norole",   null,  null, "/*", null, 0, 1),
(3, "admin",    null,  null, "/*", null, 1, 1),
(4, "jwt",      "jwt", null, "/*", null, 1, 1);

-- domain 如果同一用户有多个角色,系统可以通过访问的域名自动分配角色
INSERT INTO `zgo_role`(`id`, `kid`, `name`, `domain`, `status`, `organization`) VALUES 
(1, "admin",   "管理员", null, 1, null),
(2, "normal",  "正常",   null, 1, null),
(3, "group1",  "分组1",  null, 1, null),
(4, "invalid", "作废",   null, 0, null);

INSERT INTO `zgo_role_gateway`(`id`, `role_id`, `gateway`) VALUES 
(1, 1, "admin"),
(2, 2, "jwt");

INSERT INTO `zgo_role_role`(`id`, `owner_id`, `child_id`) VALUES 
(1, 3, 1);

INSERT INTO `zgo_user_role`(`id`, `user_id`, `role_id`) VALUES 
(1, 1, 1);

