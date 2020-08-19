-- -------------------------------------------------------
-- build by cmd/db/mysql/mysql.go
-- time: 2020-08-18 12:28:05 CST
-- -------------------------------------------------------
-- 表结构
-- -------------------------------------------------------
-- 账户实体
CREATE TABLE `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `account` varchar(255) DEFAULT NULL COMMENT '账户',
  `account_type` tinyint(4) DEFAULT '1' COMMENT '账户类型',
  `platform_id` int(11) DEFAULT NULL COMMENT '账户归属平台',
  `verify_type` tinyint(4) DEFAULT '1' COMMENT '校验方式',
  `pid` int(11) DEFAULT NULL COMMENT '上级账户',
  `password` varchar(255) DEFAULT NULL COMMENT '登录密码',
  `password_salt` varchar(255) DEFAULT NULL COMMENT '密码盐值',
  `password_type` varchar(16) DEFAULT NULL COMMENT '校验方式',
  `user_id` int(11) DEFAULT NULL COMMENT '用户标识',
  `role_id` int(11) DEFAULT NULL COMMENT '角色标识',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '账户描述',
  `oa2_token` varchar(1024) DEFAULT NULL COMMENT 'oauth2令牌',
  `oa2_expired` timestamp DEFAULT NULL COMMENT 'oauth2过期时间',
  `token_fake` varchar(1024) DEFAULT NULL COMMENT '伪造令牌',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_2` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_3` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_2` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_3` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_account(`account`,`account_type`,`platform_id`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 第三方登陆实体
CREATE TABLE `oauth2_platform` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `platform` varchar(32) NOT NULL COMMENT '平台',
  `agent_id` varchar(255) DEFAULT NULL COMMENT '代理商标识',
  `suite_id` varchar(255) DEFAULT NULL COMMENT '套件标识',
  `app_id` varchar(255) DEFAULT NULL COMMENT '应用标识',
  `app_secret` varchar(1024) DEFAULT NULL COMMENT '应用密钥',
  `authorize_url` varchar(2048) DEFAULT NULL COMMENT '认证地址',
  `token_url` varchar(2048) DEFAULT NULL COMMENT '令牌地址',
  `profile_url` varchar(2048) DEFAULT NULL COMMENT '个人资料地址',
  `domain_def` varchar(128) DEFAULT NULL COMMENT '默认域名',
  `domain_check` varchar(255) DEFAULT NULL COMMENT '域名认证',
  `js_secret` varchar(255) DEFAULT NULL COMMENT 'javascript密钥',
  `state_secret` varchar(255) DEFAULT NULL COMMENT '回调state密钥',
  `callback` tinyint(4) DEFAULT 0 COMMENT '是否支持回调',
  `cb_encrypt` tinyint(4) DEFAULT 0 COMMENT '回调是否加密',
  `cb_token` varchar(255) DEFAULT NULL COMMENT '加密令牌',
  `cb_encoding` varchar(255) DEFAULT NULL COMMENT '加密编码',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_2` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_3` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_2` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_3` int(11) DEFAULT NULL COMMENT '备用字段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 第三方登陆实体
CREATE TABLE `oauth2_token` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `oauth2_id` int(11) DEFAULT NULL COMMENT '平台',
  `access_token` varchar(1024) DEFAULT NULL COMMENT '代理商标识',
  `expires_in` int(11) DEFAULT 7200 COMMENT '有限期间隔',
  `create_time` timestamp DEFAULT NULL COMMENT '凭据创建时间',
  `sync_lock` tinyint(4) DEFAULT 0 COMMENT '同步锁',
  `call_count` int(11) DEFAULT 0 COMMENT '调用次数',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_2` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_3` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_2` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_3` int(11) DEFAULT NULL COMMENT '备用字段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 第三方授权实体
CREATE TABLE `oauth2_client` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `client_key` varchar(255) DEFAULT NULL COMMENT '客户端标识',
  `audience` varchar(1024) DEFAULT NULL COMMENT '账户接受平台',
  `issuer` varchar(1024) DEFAULT NULL COMMENT '账户签发平台',
  `expired` int(11) DEFAULT NULL COMMENT '令牌有效期',
  `token_type` varchar(32) DEFAULT NULL COMMENT '令牌类型',
  `s_method` varchar(32) DEFAULT NULL COMMENT '令牌方法',
  `s_secret` varchar(1024) DEFAULT NULL COMMENT '令牌密钥',
  `token_getter` varchar(32) DEFAULT NULL COMMENT '令牌获取方法',
  `signin_url` varchar(2048) DEFAULT NULL COMMENT '登陆地址',
  `signin_force` tinyint(4) DEFAULT 0 COMMENT '强制跳转登陆',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_2` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_3` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_2` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_3` int(11) DEFAULT NULL COMMENT '备用字段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 第三方授权实体
CREATE TABLE `oauth2_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `client_id` int(11) DEFAULT NULL COMMENT '客户端标识',
  `secret` varchar(1024) DEFAULT NULL COMMENT '密钥',
  `expired` timestamp DEFAULT NULL COMMENT '授权有效期',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_2` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_3` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_2` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_3` int(11) DEFAULT NULL COMMENT '备用字段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户实体
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `uid` varchar(64) DEFAULT NULL COMMENT '唯一标识',
  `name` varchar(64) DEFAULT NULL COMMENT '用户名',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  UNIQUE udx_user_uid(`uid`),
  UNIQUE udx_user_name(`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户详情实体
CREATE TABLE `user_detail` (
  `id` int(11) NOT NULL COMMENT '唯一标识',
  `nickname` varchar(64) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(512) DEFAULT NULL COMMENT '头像',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_2` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `string_3` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_2` int(11) DEFAULT NULL COMMENT '备用字段',
  `number_3` int(11) DEFAULT NULL COMMENT '备用字段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户消息实体
CREATE TABLE `user_message` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `uid` varchar(64) DEFAULT NULL COMMENT '索引',
  `avatar` varchar(512) DEFAULT NULL COMMENT '头像',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `datetime` timestamp DEFAULT NULL COMMENT '日期',
  `type` varchar(16) DEFAULT NULL COMMENT '类型',
  `read` tinyint(4) DEFAULT NULL COMMENT '已读',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `clickClose` tinyint(4) DEFAULT NULL COMMENT '关闭',
  `status` tinyint(4) DEFAULT NULL COMMENT '状态',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_user_message_uid(`uid`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 角色实体
CREATE TABLE `role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `uid` varchar(64) DEFAULT NULL COMMENT '唯一标识',
  `name` varchar(64) DEFAULT NULL COMMENT '角色名',
  `description` varchar(128) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `domain` varchar(255) DEFAULT NULL COMMENT '域',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_role_uid(`uid`),
  UNIQUE udx_role_name(`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 角色角色实体
CREATE TABLE `role_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `owner_id` int(11) DEFAULT NULL COMMENT '父节点标识',
  `child_id` int(11) DEFAULT NULL COMMENT '子节点标识',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户角色实体
CREATE TABLE `user_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `user_id` int(11) DEFAULT NULL COMMENT '账户标识',
  `role_id` int(11) DEFAULT NULL COMMENT '客户端标识',
  `expired` int(11) DEFAULT NULL COMMENT '授权有效期',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 资源实体
CREATE TABLE `resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `resource` varchar(64) DEFAULT NULL COMMENT '资源名',
  `domain` varchar(255) DEFAULT NULL COMMENT '域名',
  `methods` varchar(64) DEFAULT NULL COMMENT '方法',
  `path` varchar(255) DEFAULT NULL COMMENT '路径',
  `netmask` varchar(64) DEFAULT NULL COMMENT '网络掩码',
  `allow` tinyint(4) DEFAULT NULL COMMENT '允许vs拒绝',
  `description` varchar(128) DEFAULT NULL COMMENT '描述',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  INDEX idx_resource_name(`resource`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 资源角色实体
CREATE TABLE `resource_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `role_id` int(11) DEFAULT NULL COMMENT '角色',
  `resource` varchar(64) DEFAULT NULL COMMENT '资源名',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 资源角色实体
CREATE TABLE `resource_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `user_id` int(11) DEFAULT NULL COMMENT '用户',
  `resource` varchar(64) DEFAULT NULL COMMENT '资源名',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 菜单实体
CREATE TABLE `menu` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `uid` varchar(32) DEFAULT NULL COMMENT '唯一标识',
  `name` varchar(64) DEFAULT NULL COMMENT '菜单名称',
  `local` varchar(128) DEFAULT NULL COMMENT '菜单名称',
  `sequence` tinyint(4) DEFAULT 64 COMMENT '排序值',
  `group` varchar(64) DEFAULT NULL COMMENT '菜单分组',
  `group_local` varchar(64) DEFAULT NULL COMMENT '菜单分组',
  `icon` varchar(255) DEFAULT NULL COMMENT '图标',
  `router` varchar(255) DEFAULT NULL COMMENT '访问路由',
  `memo` varchar(255) DEFAULT NULL COMMENT '备注',
  `domain` varchar(255) DEFAULT NULL COMMENT '域名',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_menu_uid(`uid`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 角色自定义菜单实体
CREATE TABLE `menu_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `pid` int(11) DEFAULT NULL COMMENT '父节点',
  `uid` varchar(32) DEFAULT NULL COMMENT '唯一标识',
  `name` varchar(64) DEFAULT NULL COMMENT '菜单名称',
  `local` varchar(128) DEFAULT NULL COMMENT '菜单名称',
  `sequence` tinyint(4) DEFAULT 64 COMMENT '排序值',
  `role_id` int(11) DEFAULT NULL COMMENT '角色 ID',
  `role_uid` varchar(64) DEFAULT NULL COMMENT '角色 UID',
  `menu_id` int(11) DEFAULT NULL COMMENT '菜单 ID',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_menu_role_uid(`uid`),
  INDEX idx_role_uid(`role_uid`),
  INDEX idx_parent_id(`pid`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 角色自定义菜单实体
CREATE TABLE `menu_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `role_id` int(11) DEFAULT NULL COMMENT '角色 ID',
  `role_uid` varchar(64) DEFAULT NULL COMMENT '角色 UID',
  `user_id` int(11) DEFAULT NULL COMMENT '用户 ID',
  `user_uid` varchar(64) DEFAULT NULL COMMENT '用户 UID',
  `menu_ids` varchar(4096) DEFAULT NULL COMMENT '菜单 ID',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  INDEX idx_role_uid(`role_uid`),
  INDEX idx_user_uid(`user_uid`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 菜单事件实体
CREATE TABLE `menu_action` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `menu_id` int(11) DEFAULT NULL COMMENT '菜单 ID',
  `role_id` int(11) DEFAULT NULL COMMENT '角色 ID',
  `role_uid` varchar(64) DEFAULT NULL COMMENT '角色 UID',
  `code` varchar(64) DEFAULT NULL COMMENT '动作编号',
  `name` varchar(64) DEFAULT NULL COMMENT '动作名称',
  `disable` tinyint(4) DEFAULT 0 COMMENT '状态',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  INDEX idx_role_uid(`role_uid`),
  INDEX idx_menu_action_code(`code`),
  INDEX idx_menu_action_name(`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 通用标签实体
CREATE TABLE `tag_common` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `owner_id` int(11) DEFAULT NULL COMMENT '归属id',
  `type` tinyint(4) DEFAULT NULL COMMENT '标签类型',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_tag_common_uid(`owner_id`,`type`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 国际化实体
CREATE TABLE `i18n_language` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `mid` varchar(255) DEFAULT NULL COMMENT 'message id',
  `lang` varchar(16) DEFAULT NULL COMMENT '语言',
  `description` varchar(64) DEFAULT NULL COMMENT '描述',
  `left_delim` varchar(16) DEFAULT NULL COMMENT '定界符',
  `right_delim` varchar(16) DEFAULT NULL COMMENT '定界符',
  `zero` varchar(255) DEFAULT NULL COMMENT 'zero',
  `one` varchar(255) DEFAULT NULL COMMENT 'one',
  `two` varchar(255) DEFAULT NULL COMMENT 'two',
  `few` varchar(255) DEFAULT NULL COMMENT 'few',
  `many` varchar(255) DEFAULT NULL COMMENT 'many',
  `other` varchar(255) DEFAULT NULL COMMENT 'other',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_i18n_message_id(`mid`,`lang`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------

-- -------------------------------------------------------
-- 表外键
-- -------------------------------------------------------
ALTER TABLE `account`
ADD CONSTRAINT `fk_platform_id` FOREIGN KEY (`platform_id`)  REFERENCES `oauth2_platform` (`id`),
ADD CONSTRAINT `fk_account_user` FOREIGN KEY (`user_id`)  REFERENCES `user` (`id`),
ADD CONSTRAINT `fk_account_role` FOREIGN KEY (`role_id`)  REFERENCES `role` (`id`);

ALTER TABLE `oauth2_token`
ADD CONSTRAINT `fk_oa2_token_id` FOREIGN KEY (`oauth2_id`)  REFERENCES `oauth2_platform` (`id`);

ALTER TABLE `oauth2_account`
ADD CONSTRAINT `fk_oa2_client_id` FOREIGN KEY (`client_id`)  REFERENCES `oauth2_client` (`id`);

ALTER TABLE `user_detail`
ADD CONSTRAINT `fk_user_detail` FOREIGN KEY (`id`)  REFERENCES `user` (`id`);

ALTER TABLE `role_role`
ADD CONSTRAINT `fk_role_owner_id` FOREIGN KEY (`owner_id`)  REFERENCES `role` (`id`),
ADD CONSTRAINT `fk_role_child_id` FOREIGN KEY (`child_id`)  REFERENCES `role` (`id`);

ALTER TABLE `user_role`
ADD CONSTRAINT `fk_role_role_id` FOREIGN KEY (`role_id`)  REFERENCES `role` (`id`),
ADD CONSTRAINT `fk_role_user_id` FOREIGN KEY (`user_id`)  REFERENCES `user` (`id`);

ALTER TABLE `resource_role`
ADD CONSTRAINT `fk_resource_role_id` FOREIGN KEY (`role_id`)  REFERENCES `role` (`id`),
ADD CONSTRAINT `fk_resource_role_res` FOREIGN KEY (`resource`)  REFERENCES `resource` (`resource`);

ALTER TABLE `resource_user`
ADD CONSTRAINT `fk_resource_user_res` FOREIGN KEY (`resource`)  REFERENCES `resource` (`resource`),
ADD CONSTRAINT `fk_resource_user_id` FOREIGN KEY (`user_id`)  REFERENCES `user` (`id`);

ALTER TABLE `menu_role`
ADD CONSTRAINT `fk_menu_role_role_id` FOREIGN KEY (`role_id`)  REFERENCES `role` (`id`),
ADD CONSTRAINT `fk_menu_role_menu_id` FOREIGN KEY (`menu_id`)  REFERENCES `menu` (`id`);

ALTER TABLE `menu_user`
ADD CONSTRAINT `fk_menu_user_role_id` FOREIGN KEY (`role_id`)  REFERENCES `role` (`id`),
ADD CONSTRAINT `fk_menu_user_user_id` FOREIGN KEY (`user_id`)  REFERENCES `user` (`id`);

ALTER TABLE `menu_action`
ADD CONSTRAINT `fk_menu_action_menu_id` FOREIGN KEY (`menu_id`)  REFERENCES `menu` (`id`),
ADD CONSTRAINT `fk_menu_action_role_id` FOREIGN KEY (`role_id`)  REFERENCES `role` (`id`);

-- -------------------------------------------------------
-- -------------------------------------------------------
-- insert into 
-- -------------------------------------------------------
-- INSERT INTO `account`(`id`, `account`, `account_type`, `platform_id`, `verify_type`, `pid`, `password`, `password_salt`, `password_type`, `user_id`, `role_id`, `status`, `description`, `oa2_token`, `oa2_expired`, `token_fake`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `string_2`, `string_3`, `number_1`, `number_2`, `number_3`) VALUES ();
-- INSERT INTO `oauth2_platform`(`id`, `platform`, `agent_id`, `suite_id`, `app_id`, `app_secret`, `authorize_url`, `token_url`, `profile_url`, `domain_def`, `domain_check`, `js_secret`, `state_secret`, `callback`, `cb_encrypt`, `cb_token`, `cb_encoding`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `string_2`, `string_3`, `number_1`, `number_2`, `number_3`) VALUES ();
-- INSERT INTO `oauth2_token`(`id`, `oauth2_id`, `access_token`, `expires_in`, `create_time`, `sync_lock`, `call_count`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `string_2`, `string_3`, `number_1`, `number_2`, `number_3`) VALUES ();
-- INSERT INTO `oauth2_client`(`id`, `client_key`, `audience`, `issuer`, `expired`, `token_type`, `s_method`, `s_secret`, `token_getter`, `signin_url`, `signin_force`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `string_2`, `string_3`, `number_1`, `number_2`, `number_3`) VALUES ();
-- INSERT INTO `oauth2_account`(`id`, `client_id`, `secret`, `expired`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `string_2`, `string_3`, `number_1`, `number_2`, `number_3`) VALUES ();
-- INSERT INTO `user`(`id`, `uid`, `name`, `status`) VALUES ();
-- INSERT INTO `user_detail`(`id`, `nickname`, `avatar`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `string_2`, `string_3`, `number_1`, `number_2`, `number_3`) VALUES ();
-- INSERT INTO `user_message`(`id`, `uid`, `avatar`, `title`, `datetime`, `type`, `read`, `description`, `clickClose`, `status`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `role`(`id`, `uid`, `name`, `description`, `status`, `domain`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `role_role`(`id`, `owner_id`, `child_id`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `user_role`(`id`, `user_id`, `role_id`, `expired`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `resource`(`id`, `resource`, `domain`, `methods`, `path`, `netmask`, `allow`, `description`, `status`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `resource_role`(`id`, `role_id`, `resource`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `resource_user`(`id`, `user_id`, `resource`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `menu`(`id`, `uid`, `name`, `local`, `sequence`, `group`, `group_local`, `icon`, `router`, `memo`, `domain`, `status`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `menu_role`(`id`, `pid`, `uid`, `name`, `local`, `sequence`, `role_id`, `role_uid`, `menu_id`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `menu_user`(`id`, `role_id`, `role_uid`, `user_id`, `user_uid`, `menu_ids`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `menu_action`(`id`, `menu_id`, `role_id`, `role_uid`, `code`, `name`, `disable`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `tag_common`(`id`, `owner_id`, `type`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `i18n_language`(`id`, `mid`, `lang`, `description`, `left_delim`, `right_delim`, `zero`, `one`, `two`, `few`, `many`, `other`, `status`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();

-- -------------------------------------------------------
-- -------------------------------------------------------
-- drop table 
-- -------------------------------------------------------
-- ALTER TABLE `account`
-- DROP FOREIGN KEY `fk_platform_id`,
-- DROP FOREIGN KEY `fk_account_user`,
-- DROP FOREIGN KEY `fk_account_role`;
-- ALTER TABLE `oauth2_token`
-- DROP FOREIGN KEY `fk_oa2_token_id`;
-- ALTER TABLE `oauth2_account`
-- DROP FOREIGN KEY `fk_oa2_client_id`;
-- ALTER TABLE `user_detail`
-- DROP FOREIGN KEY `fk_user_detail`;
-- ALTER TABLE `role_role`
-- DROP FOREIGN KEY `fk_role_owner_id`,
-- DROP FOREIGN KEY `fk_role_child_id`;
-- ALTER TABLE `user_role`
-- DROP FOREIGN KEY `fk_role_user_id`,
-- DROP FOREIGN KEY `fk_role_role_id`;
-- ALTER TABLE `resource_role`
-- DROP FOREIGN KEY `fk_resource_role_id`,
-- DROP FOREIGN KEY `fk_resource_role_res`;
-- ALTER TABLE `resource_user`
-- DROP FOREIGN KEY `fk_resource_user_id`,
-- DROP FOREIGN KEY `fk_resource_user_res`;
-- ALTER TABLE `menu_role`
-- DROP FOREIGN KEY `fk_menu_role_role_id`,
-- DROP FOREIGN KEY `fk_menu_role_menu_id`;
-- ALTER TABLE `menu_user`
-- DROP FOREIGN KEY `fk_menu_user_user_id`,
-- DROP FOREIGN KEY `fk_menu_user_role_id`;
-- ALTER TABLE `menu_action`
-- DROP FOREIGN KEY `fk_menu_action_menu_id`,
-- DROP FOREIGN KEY `fk_menu_action_role_id`;
-- 
-- DROP TABLE IF EXISTS `account`;
-- DROP TABLE IF EXISTS `oauth2_platform`;
-- DROP TABLE IF EXISTS `oauth2_token`;
-- DROP TABLE IF EXISTS `oauth2_client`;
-- DROP TABLE IF EXISTS `oauth2_account`;
-- DROP TABLE IF EXISTS `user`;
-- DROP TABLE IF EXISTS `user_detail`;
-- DROP TABLE IF EXISTS `user_message`;
-- DROP TABLE IF EXISTS `role`;
-- DROP TABLE IF EXISTS `role_role`;
-- DROP TABLE IF EXISTS `user_role`;
-- DROP TABLE IF EXISTS `resource`;
-- DROP TABLE IF EXISTS `resource_role`;
-- DROP TABLE IF EXISTS `resource_user`;
-- DROP TABLE IF EXISTS `menu`;
-- DROP TABLE IF EXISTS `menu_role`;
-- DROP TABLE IF EXISTS `menu_user`;
-- DROP TABLE IF EXISTS `menu_action`;
-- DROP TABLE IF EXISTS `tag_common`;
-- DROP TABLE IF EXISTS `i18n_language`;

-- -------------------------------------------------------