-- -------------------------------------------------------
-- build by cmd/db/mysql/mysql.go
-- time: 2021-01-14 10:03:28 CST
-- -------------------------------------------------------
-- 表结构
-- -------------------------------------------------------
-- 账户实体
CREATE TABLE `zgo_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `pid` int(11) DEFAULT NULL COMMENT '上级账户',
  `account` varchar(255) NOT NULL COMMENT '账户',
  `account_type` tinyint(4) DEFAULT '1' COMMENT '账户类型',
  `platform_kid` varchar(64) DEFAULT NULL COMMENT '账户归属平台',
  `user_id` int(11) NOT NULL COMMENT '用户标识',
  `role_id` int(11) DEFAULT NULL COMMENT '角色标识',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `password` varchar(255) DEFAULT NULL COMMENT '登录密码',
  `password_salt` varchar(255) DEFAULT NULL COMMENT '密码盐值',
  `password_type` varchar(16) DEFAULT NULL COMMENT '密码方式',
  `verify_secret` varchar(255) DEFAULT NULL COMMENT '校验密钥',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '账户描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_account(`account`,`account_type`,`platform_kid`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户实体
CREATE TABLE `zgo_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `pid` int(11) DEFAULT NULL COMMENT '用户归属',
  `kid` varchar(64) NOT NULL COMMENT '唯一标识',
  `type` varchar(64) NOT NULL DEFAULT 'usr' COMMENT '账户类型',
  `name` varchar(64) NOT NULL COMMENT '用户名',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `delete` tinyint(4) DEFAULT 0 COMMENT '删除标识',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_user_kid(`kid`),
  INDEX idx_user_type(`type`),
  INDEX idx_user_name(`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户详情实体
CREATE TABLE `zgo_user_detail` (
  `id` int(11) NOT NULL COMMENT '唯一标识',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `description` varchar(255) DEFAULT NULL COMMENT '个人描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户安全实体
CREATE TABLE `zgo_user_security` (
  `id` int(11) NOT NULL COMMENT '唯一标识',
  `mfa13` varchar(1024) DEFAULT NULL COMMENT 'mfa密钥',
  `phone` varchar(32) DEFAULT NULL COMMENT '密保电话',
  `email` varchar(128) DEFAULT NULL COMMENT '备用邮箱',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 人员实体
CREATE TABLE `zgo_person` (
  `id` int(11) NOT NULL COMMENT '唯一标识',
  `unique_name` varchar(64) NOT NULL COMMENT '唯一用户名',
  `first_name` varchar(32) DEFAULT NULL COMMENT '用户名',
  `last_name` varchar(32) DEFAULT NULL COMMENT '用户名',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户ID',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_person_unique_name(`unique_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 租户实体
CREATE TABLE `zgo_organization` (
  `id` int(11) NOT NULL COMMENT '唯一标识',
  `code` varchar(64) NOT NULL COMMENT '租户标识',
  `hosted` tinyint(4) DEFAULT 0 COMMENT '租户被托管',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_organization_code(`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 租户用户主键
CREATE TABLE `zgo_organization_user` (
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `org_cod` varchar(64) NOT NULL COMMENT '租户ID',
  `union_kid` varchar(64) NOT NULL COMMENT '唯一标识',
  `name` varchar(64) NOT NULL COMMENT '用户名',
  `custom_id` varchar(255) DEFAULT NULL COMMENT '唯一标识',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_org_user_kid(`user_id`,`org_cod`),
  UNIQUE udx_org_user_union_kid(`org_cod`,`union_kid`),
  UNIQUE udx_org_user_name(`org_cod`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 人归一方式
CREATE TABLE `zgo_user_union` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `pid` int(11) DEFAULT NULL COMMENT '唯一标识',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `type` varchar(64)  NOT NULL COMMENT '归一方式',
  `type_id` varchar(255) DEFAULT NULL COMMENT '归一标识',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  INDEX idx_user_union_type(`type`),
  INDEX idx_user_union_type_id(`type_id`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 策略服务实体
CREATE TABLE `zgo_app_service` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `name` varchar(64) DEFAULT NULL COMMENT '服务名称',
  `code` varchar(64) DEFAULT NULL COMMENT '服务编码',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_app_service_name(`name`),
  UNIQUE udx_app_service_code(`code`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 策略服务实体
CREATE TABLE `zgo_app_service_org` (
  `svc_id` int(11) NOT NULL COMMENT '服务标识',
  `org_cod` varchar(64) NOT NULL COMMENT '租户标识',
  `expired` timestamp DEFAULT NULL COMMENT '授权有效期',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 策略服务受众实体
CREATE TABLE `zgo_app_service_audience` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `svc_id` int(11) NOT NULL COMMENT '服务标识',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `audience` varchar(255) DEFAULT NULL COMMENT '受众域',
  `resource` varchar(255) DEFAULT NULL COMMENT '受众源',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_app_service_audience(`audience`,`resource`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 角色实体
CREATE TABLE `zgo_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `kid` varchar(64) DEFAULT NULL COMMENT '唯一标识',
  `svc_id` int(11) DEFAULT NULL COMMENT '服务标识',
  `name` varchar(64) DEFAULT NULL COMMENT '角色名称',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `org_adm` tinyint(4) DEFAULT 1 COMMENT '管理员',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_role_kid(`kid`),
  UNIQUE udx_role_name(`name`,`org_cod`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 角色角色实体
CREATE TABLE `zgo_role_role` (
  `pid` int(11) DEFAULT NULL COMMENT '父节点标识',
  `cid` int(11) DEFAULT NULL COMMENT '子节点标识',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户角色实体
CREATE TABLE `zgo_user_role` (
  `user_id` int(11) DEFAULT NULL COMMENT '账户标识',
  `role_id` int(11) DEFAULT NULL COMMENT '角色标识',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `expired` timestamp DEFAULT NULL COMMENT '授权有效期',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户角色实体
CREATE TABLE `zgo_role_policy` (
  `role_id` int(11) DEFAULT NULL COMMENT '角色标识',
  `plcy_id` int(11) DEFAULT NULL COMMENT '策略标识',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 用户角色实体
CREATE TABLE `zgo_user_policy` (
  `user_id` int(11) DEFAULT NULL COMMENT '账户标识',
  `plcy_id` int(11) DEFAULT NULL COMMENT '策略标识',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `expired` timestamp DEFAULT NULL COMMENT '授权有效期',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 策略服务操作实体
CREATE TABLE `zgo_policy_service_action` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `svc_id` int(11) DEFAULT NULL COMMENT '服务标识',
  `name` varchar(64) DEFAULT NULL COMMENT '操作名称',
  `resource` varchar(255) DEFAULT NULL COMMENT '资源列表',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_policy_service_action_name(`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 策略实体
CREATE TABLE `zgo_policy` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `ver` varchar(16) DEFAULT '1' COMMENT '版本',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '组织标识',
  `name` varchar(64) DEFAULT NULL COMMENT '策略名称',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_policy_name(`org_cod`,`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 策略声明
CREATE TABLE `zgo_policy_statement` (
  `plcy_id` int(11) DEFAULT NULL COMMENT '策略标识',
  `effect` tinyint(4) DEFAULT NULL COMMENT '相应',
  `action` varchar(255) DEFAULT NULL COMMENT '操作',
  `resource` varchar(255) DEFAULT NULL COMMENT '资源',
  `condition` varchar(255) DEFAULT NULL COMMENT '条件',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 第三方登陆实体
CREATE TABLE `zgo_platform` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
  `kid` varchar(64)  NOT NULL COMMENT '三方标识',
  `type` varchar(32) NOT NULL COMMENT '平台标识',
  `signin` tinyint(4) DEFAULT 0 COMMENT '可登陆',
  `org_cod` varchar(64) DEFAULT NULL COMMENT '租户标识',
  `avatar` varchar(255) DEFAULT NULL COMMENT '平台头像',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '平台描述',
  `app_id` varchar(128) DEFAULT NULL COMMENT '应用标识',
  `app_secret` varchar(1024) DEFAULT NULL COMMENT '应用密钥',
  `agent_id` varchar(128) DEFAULT NULL COMMENT '代理商标识',
  `agent_secret` varchar(1024) DEFAULT NULL COMMENT '代理商密钥',
  `suite_id` varchar(128) DEFAULT NULL COMMENT '套件标识',
  `suite_secret` varchar(1024) DEFAULT NULL COMMENT '套件密钥',
  `authorize_url` varchar(1024) DEFAULT NULL COMMENT '认证地址',
  `token_url` varchar(1024) DEFAULT NULL COMMENT '令牌地址',
  `profile_url` varchar(1024) DEFAULT NULL COMMENT '个人资料地址',
  `signin_url` varchar(128) DEFAULT NULL COMMENT '重定向回调地址',
  `js_secret` varchar(255) DEFAULT NULL COMMENT 'javascript密钥',
  `state_secret` varchar(255) DEFAULT NULL COMMENT '回调state密钥',
  `callback` tinyint(4) DEFAULT 0 COMMENT '是否支持回调',
  `cb_domain` varchar(128) DEFAULT NULL COMMENT '默认域名',
  `cb_scheme` varchar(16) DEFAULT 'https' COMMENT '默认协议',
  `cb_encrypt` tinyint(4) DEFAULT 0 COMMENT '回调是否加密',
  `cb_token` varchar(255) DEFAULT NULL COMMENT '加密令牌',
  `cb_encoding` varchar(255) DEFAULT NULL COMMENT '加密编码',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_oauth2_platform_kid(`kid`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 通信令牌实体
CREATE TABLE `zgo_token` (
  `token_kid` varchar(64) DEFAULT NULL COMMENT '角色标识',
  `account_id` int(11) DEFAULT NULL COMMENT '令牌归属',
  `access_token` varchar(4096) DEFAULT NULL COMMENT '访问令牌',
  `access_expires` int(11) DEFAULT 7200 COMMENT '访问令牌',
  `refresh_token` varchar(255) DEFAULT NULL COMMENT '刷新令牌',
  `refresh_expires` int(11) DEFAULT 604800 COMMENT '刷新令牌',
  `call_count` int(11) DEFAULT 0 COMMENT '调用次数',
  `sync_lock` int(11) DEFAULT NULL COMMENT '同步锁',
  `refresh_count` int(11) DEFAULT 0 COMMENT '刷新次数',
  `last_ip` varchar(64) DEFAULT NULL COMMENT '上次登陆IP',
  `last_at` timestamp DEFAULT NULL COMMENT '上次登陆时间',
  `error_code` varchar(64) DEFAULT NULL COMMENT '异常类型',
  `error_message` varchar(255) DEFAULT NULL COMMENT '异常信息',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  `string_1` varchar(255) DEFAULT NULL COMMENT '备用字段',
  `number_1` int(11) DEFAULT NULL COMMENT '备用字段',
  UNIQUE udx_token_kid(`token_kid`),
  INDEX idx_token_account_id(`account_id`),
  INDEX idx_token_refresh_token(`refresh_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 系统标签实体
CREATE TABLE `zgo_tag_system` (
  `type` varchar(64) DEFAULT NULL COMMENT '标签类型',
  `bid` int(11) DEFAULT NULL COMMENT '归属id',
  `tag` varchar(64) DEFAULT NULL COMMENT '标签',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_tag_system(`type`,`bid`,`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 通用标签实体
CREATE TABLE `zgo_tag_custom` (
  `type` varchar(64) DEFAULT NULL COMMENT '标签类型',
  `bid` int(11) DEFAULT NULL COMMENT '归属id',
  `tag` varchar(64) DEFAULT NULL COMMENT '标签',
  `deleted` tinyint(4) DEFAULT 0 COMMENT '删除标识',
  `creator` varchar(64) DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT 0 COMMENT '数据版本',
  UNIQUE udx_tag_custom(`type`,`bid`,`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------
-- 国际化实体
CREATE TABLE `zgo_i18n_language` (
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
  UNIQUE udx_i18n_message_id(`mid`,`lang`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- -------------------------------------------------------

-- -------------------------------------------------------
-- 表外键
-- -------------------------------------------------------
ALTER TABLE `zgo_account`
ADD CONSTRAINT `fk_account_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`),
ADD CONSTRAINT `fk_account_pid` FOREIGN KEY (`pid`)  REFERENCES `zgo_account` (`id`),
ADD CONSTRAINT `fk_account_platform_kid` FOREIGN KEY (`platform_kid`)  REFERENCES `zgo_platform` (`kid`),
ADD CONSTRAINT `fk_account_user_id` FOREIGN KEY (`user_id`)  REFERENCES `zgo_user` (`id`),
ADD CONSTRAINT `fk_account_role_id` FOREIGN KEY (`role_id`)  REFERENCES `zgo_user_role` (`id`);

ALTER TABLE `zgo_user`
ADD CONSTRAINT `fk_user_pid` FOREIGN KEY (`pid`)  REFERENCES `zgo_user` (`id`);

ALTER TABLE `zgo_user_detail`
ADD CONSTRAINT `fk_user_detail_id` FOREIGN KEY (`id`)  REFERENCES `zgo_user` (`id`);

ALTER TABLE `zgo_user_security`
ADD CONSTRAINT `fk_user_security` FOREIGN KEY (`id`)  REFERENCES `zgo_user` (`id`);

ALTER TABLE `zgo_person`
ADD CONSTRAINT `fk_person_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`),
ADD CONSTRAINT `fk_person_id` FOREIGN KEY (`id`)  REFERENCES `zgo_user` (`id`);

ALTER TABLE `zgo_organization`
ADD CONSTRAINT `fk_organization_id` FOREIGN KEY (`id`)  REFERENCES `zgo_user` (`id`);

ALTER TABLE `zgo_organization_user`
ADD CONSTRAINT `fk_org_user_uid` FOREIGN KEY (`user_id`)  REFERENCES `zgo_user` (`id`),
ADD CONSTRAINT `fk_org_user_ocd` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`);

ALTER TABLE `zgo_user_union`
ADD CONSTRAINT `fk_user_union_pid` FOREIGN KEY (`pid`)  REFERENCES `zgo_user_union` (`id`),
ADD CONSTRAINT `fk_user_union_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`),
ADD CONSTRAINT `fk_user_union_user_id` FOREIGN KEY (`user_id`)  REFERENCES `zgo_user` (`id`);

ALTER TABLE `zgo_app_service_org`
ADD CONSTRAINT `fk_app_service_org_sid` FOREIGN KEY (`svc_id`)  REFERENCES `zgo_app_service` (`id`),
ADD CONSTRAINT `fk_app_service_org_orcd` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`);

ALTER TABLE `zgo_app_service_audience`
ADD CONSTRAINT `fk_app_service_audience_ocd` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`),
ADD CONSTRAINT `fk_app_service_audience_sid` FOREIGN KEY (`svc_id`)  REFERENCES `zgo_app_service` (`id`);

ALTER TABLE `zgo_role`
ADD CONSTRAINT `fk_role_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`),
ADD CONSTRAINT `fk_role_sid` FOREIGN KEY (`svc_id`)  REFERENCES `zgo_app_service` (`id`);

ALTER TABLE `zgo_role_role`
ADD CONSTRAINT `fk_role_role_pid` FOREIGN KEY (`pid`)  REFERENCES `zgo_role` (`id`),
ADD CONSTRAINT `fk_role_role_cid` FOREIGN KEY (`cid`)  REFERENCES `zgo_role` (`id`),
ADD CONSTRAINT `fk_role_role_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`);

ALTER TABLE `zgo_user_role`
ADD CONSTRAINT `fk_user_role_user_id` FOREIGN KEY (`user_id`)  REFERENCES `zgo_user` (`id`),
ADD CONSTRAINT `fk_user_role_role_id` FOREIGN KEY (`role_id`)  REFERENCES `zgo_role` (`id`),
ADD CONSTRAINT `fk_user_role_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`);

ALTER TABLE `zgo_role_policy`
ADD CONSTRAINT `fk_role_policy_role_id` FOREIGN KEY (`role_id`)  REFERENCES `zgo_role` (`id`),
ADD CONSTRAINT `fk_role_policy_plcy_id` FOREIGN KEY (`plcy_id`)  REFERENCES `zgo_role` (`id`),
ADD CONSTRAINT `fk_role_policy_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`);

ALTER TABLE `zgo_user_policy`
ADD CONSTRAINT `fk_role_role_user_id` FOREIGN KEY (`user_id`)  REFERENCES `zgo_user` (`id`),
ADD CONSTRAINT `fk_role_role_plcy_id` FOREIGN KEY (`plcy_id`)  REFERENCES `zgo_role` (`id`),
ADD CONSTRAINT `fk_user_policy_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`);

ALTER TABLE `zgo_policy_service_action`
ADD CONSTRAINT `fk_policy_service_action_sid` FOREIGN KEY (`svc_id`)  REFERENCES `zgo_app_service` (`id`);

ALTER TABLE `zgo_policy`
ADD CONSTRAINT `fk_policy_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`code`);

ALTER TABLE `zgo_policy_statement`
ADD CONSTRAINT `fk_policy_statement_plcy_id` FOREIGN KEY (`plcy_id`)  REFERENCES `zgo_policy` (`id`);

ALTER TABLE `zgo_platform`
ADD CONSTRAINT `fk_platform_org_cod` FOREIGN KEY (`org_cod`)  REFERENCES `zgo_organization` (`cod`);

-- -------------------------------------------------------
-- -------------------------------------------------------
-- insert into 
-- -------------------------------------------------------
-- INSERT INTO `zgo_account`(`id`, `pid`, `account`, `account_type`, `platform_kid`, `user_id`, `role_id`, `org_cod`, `password`, `password_salt`, `password_type`, `verify_secret`, `status`, `description`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_user`(`id`, `pid`, `kid`, `type`, `name`, `status`, `delete`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_user_detail`(`id`, `avatar`, `description`, `creator`, `created_at`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_user_security`(`id`, `mfa13`, `phone`, `email`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_person`(`id`, `unique_name`, `first_name`, `last_name`, `org_cod`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_organization`(`id`, `code`, `hosted`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_organization_user`(`user_id`, `org_cod`, `union_kid`, `name`, `custom_id`, `status`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_user_union`(`id`, `pid`, `org_cod`, `user_id`, `type`, `type_id`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_app_service`(`id`, `name`, `code`, `status`, `description`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_app_service_org`(`svc_id`, `org_cod`, `expired`, `status`, `description`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_app_service_audience`(`id`, `svc_id`, `org_cod`, `audience`, `resource`, `status`, `description`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_role`(`id`, `kid`, `svc_id`, `name`, `org_cod`, `org_adm`, `status`, `description`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_role_role`(`pid`, `cid`, `org_cod`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_user_role`(`user_id`, `role_id`, `org_cod`, `expired`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_role_policy`(`role_id`, `plcy_id`, `org_cod`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_user_policy`(`user_id`, `plcy_id`, `org_cod`, `expired`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_policy_service_action`(`id`, `svc_id`, `name`, `resource`, `status`, `description`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_policy`(`id`, `ver`, `org_cod`, `name`, `status`, `description`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_policy_statement`(`plcy_id`, `effect`, `action`, `resource`, `condition`, `description`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_platform`(`id`, `kid`, `type`, `signin`, `org_cod`, `avatar`, `status`, `description`, `app_id`, `app_secret`, `agent_id`, `agent_secret`, `suite_id`, `suite_secret`, `authorize_url`, `token_url`, `profile_url`, `signin_url`, `js_secret`, `state_secret`, `callback`, `cb_domain`, `cb_scheme`, `cb_encrypt`, `cb_token`, `cb_encoding`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_token`(`token_kid`, `account_id`, `access_token`, `access_expires`, `refresh_token`, `refresh_expires`, `call_count`, `sync_lock`, `refresh_count`, `last_ip`, `last_at`, `error_code`, `error_message`, `creator`, `created_at`, `updated_at`, `version`, `string_1`, `number_1`) VALUES ();
-- INSERT INTO `zgo_tag_system`(`type`, `bid`, `tag`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_tag_custom`(`type`, `bid`, `tag`, `deleted`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();
-- INSERT INTO `zgo_i18n_language`(`mid`, `lang`, `description`, `left_delim`, `right_delim`, `zero`, `one`, `two`, `few`, `many`, `other`, `status`, `creator`, `created_at`, `updated_at`, `version`) VALUES ();

-- -------------------------------------------------------
-- -------------------------------------------------------
-- drop table 
-- -------------------------------------------------------
-- ALTER TABLE `zgo_account`
-- DROP FOREIGN KEY `fk_account_org_cod`,
-- DROP FOREIGN KEY `fk_account_pid`,
-- DROP FOREIGN KEY `fk_account_platform_kid`,
-- DROP FOREIGN KEY `fk_account_user_id`,
-- DROP FOREIGN KEY `fk_account_role_id`;
-- ALTER TABLE `zgo_user`
-- DROP FOREIGN KEY `fk_user_pid`;
-- ALTER TABLE `zgo_user_detail`
-- DROP FOREIGN KEY `fk_user_detail_id`;
-- ALTER TABLE `zgo_user_security`
-- DROP FOREIGN KEY `fk_user_security`;
-- ALTER TABLE `zgo_person`
-- DROP FOREIGN KEY `fk_person_id`,
-- DROP FOREIGN KEY `fk_person_org_cod`;
-- ALTER TABLE `zgo_organization`
-- DROP FOREIGN KEY `fk_organization_id`;
-- ALTER TABLE `zgo_organization_user`
-- DROP FOREIGN KEY `fk_org_user_ocd`,
-- DROP FOREIGN KEY `fk_org_user_uid`;
-- ALTER TABLE `zgo_user_union`
-- DROP FOREIGN KEY `fk_user_union_pid`,
-- DROP FOREIGN KEY `fk_user_union_org_cod`,
-- DROP FOREIGN KEY `fk_user_union_user_id`;
-- ALTER TABLE `zgo_app_service_org`
-- DROP FOREIGN KEY `fk_app_service_org_sid`,
-- DROP FOREIGN KEY `fk_app_service_org_orcd`;
-- ALTER TABLE `zgo_app_service_audience`
-- DROP FOREIGN KEY `fk_app_service_audience_ocd`,
-- DROP FOREIGN KEY `fk_app_service_audience_sid`;
-- ALTER TABLE `zgo_role`
-- DROP FOREIGN KEY `fk_role_sid`,
-- DROP FOREIGN KEY `fk_role_org_cod`;
-- ALTER TABLE `zgo_role_role`
-- DROP FOREIGN KEY `fk_role_role_pid`,
-- DROP FOREIGN KEY `fk_role_role_cid`,
-- DROP FOREIGN KEY `fk_role_role_org_cod`;
-- ALTER TABLE `zgo_user_role`
-- DROP FOREIGN KEY `fk_user_role_role_id`,
-- DROP FOREIGN KEY `fk_user_role_org_cod`,
-- DROP FOREIGN KEY `fk_user_role_user_id`;
-- ALTER TABLE `zgo_role_policy`
-- DROP FOREIGN KEY `fk_role_policy_role_id`,
-- DROP FOREIGN KEY `fk_role_policy_plcy_id`,
-- DROP FOREIGN KEY `fk_role_policy_org_cod`;
-- ALTER TABLE `zgo_user_policy`
-- DROP FOREIGN KEY `fk_role_role_user_id`,
-- DROP FOREIGN KEY `fk_role_role_plcy_id`,
-- DROP FOREIGN KEY `fk_user_policy_org_cod`;
-- ALTER TABLE `zgo_policy_service_action`
-- DROP FOREIGN KEY `fk_policy_service_action_sid`;
-- ALTER TABLE `zgo_policy`
-- DROP FOREIGN KEY `fk_policy_org_cod`;
-- ALTER TABLE `zgo_policy_statement`
-- DROP FOREIGN KEY `fk_policy_statement_plcy_id`;
-- ALTER TABLE `zgo_platform`
-- DROP FOREIGN KEY `fk_platform_org_cod`;
-- 
-- DROP TABLE IF EXISTS `zgo_account`;
-- DROP TABLE IF EXISTS `zgo_user`;
-- DROP TABLE IF EXISTS `zgo_user_detail`;
-- DROP TABLE IF EXISTS `zgo_user_security`;
-- DROP TABLE IF EXISTS `zgo_person`;
-- DROP TABLE IF EXISTS `zgo_organization`;
-- DROP TABLE IF EXISTS `zgo_organization_user`;
-- DROP TABLE IF EXISTS `zgo_user_union`;
-- DROP TABLE IF EXISTS `zgo_app_service`;
-- DROP TABLE IF EXISTS `zgo_app_service_org`;
-- DROP TABLE IF EXISTS `zgo_app_service_audience`;
-- DROP TABLE IF EXISTS `zgo_role`;
-- DROP TABLE IF EXISTS `zgo_role_role`;
-- DROP TABLE IF EXISTS `zgo_user_role`;
-- DROP TABLE IF EXISTS `zgo_role_policy`;
-- DROP TABLE IF EXISTS `zgo_user_policy`;
-- DROP TABLE IF EXISTS `zgo_policy_service_action`;
-- DROP TABLE IF EXISTS `zgo_policy`;
-- DROP TABLE IF EXISTS `zgo_policy_statement`;
-- DROP TABLE IF EXISTS `zgo_platform`;
-- DROP TABLE IF EXISTS `zgo_token`;
-- DROP TABLE IF EXISTS `zgo_tag_system`;
-- DROP TABLE IF EXISTS `zgo_tag_custom`;
-- DROP TABLE IF EXISTS `zgo_i18n_language`;

-- -------------------------------------------------------