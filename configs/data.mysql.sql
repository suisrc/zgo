INSERT INTO `user`(`id`, `uid`, `name`) VALUES
(1, "uid01", "admin"),
(2, "uid02", "user");

INSERT INTO `account`(`id`, `account`, `password`, `user_id`, `role_id`, `oauth2_id`) VALUES 
(1, "admin", "123456", 1, null, null),
(2, "user",  "123456", 1, null, null);

