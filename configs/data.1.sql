-- 用户 user_id=1，为超级管理员，所有的权限认证将被跳过
-- 系统初始化完成后， 最好将该角色权限禁用
-- 编码说明
-- 用户(48)： u<助记符3位><时间编码8位><ID编码8位><机器码4位><随机码24位>
-- 租户(32)： o<助记符3位><时间编码8位><ID编码8位><机器码4位><随机码8位>
-- 应用(24)： a<助记符3位><时间编码8位><机器码4位><随机码8位>
INSERT INTO `zgo_user`(`id`, `type`, `kid`, `name`, `status`) VALUES
(1,      "org", "00000000000000000000000000000000",                 "P6M-ADM", 0),
(10007,  "usr", "u00020210112xx0100071111x000000000x000000000x007", "P6M-T7",  1),
(100001, "org", "o00020210112xx1000011111000000a1",                 "赢迪",    1),
(100002, "usr", "u00020210112xx1000021111x000000000x000000000x002", "USR-IT2", 1),
(100003, "usr", "u00020210112xx1000031111x000000000x000000000x003", "USR-IT3", 1);

-- zgo_user自增从 100000001 开始, 小于100000001为保留字段， 用于人工录入修正
ALTER TABLE `zgo_user` AUTO_INCREMENT = 100000001;


-- 租户 code=P6M， 为平台租户
INSERT INTO `zgo_organization`(`id`, `code`) VALUES
(1,      "P6M"),
(100001, "ORGCM3558");

-- 人员
INSERT INTO `zgo_person`(`id`, `unique_name`, `first_name`, `last_name`) VALUES
(10007,  "P6M-测试7", null, null),
(100002, "USR-测试2", null, null),
(100003, "USR-测试3", null, null);

-- 租户用户
-- 租户编码(40)： u<助记符3位><租户ID编码8位><时间编码8位><用户ID编码8位><机器码4位><随机码8位>
INSERT INTO `zgo_organization_user`(`user_id`, `org_cod`, `ouid`, `name`) VALUES
(100002, "ORGCM3558", "u000xx10000120210112xx100002111100000001", "赢迪测试2"),
(100003, "ORGCM3558", "u000xx10000120210112xx100003111100000001", "赢迪测试3");

-- 账户
INSERT INTO `zgo_account`(`id`, `pid`, `user_id`, `account`, `account_type`, `password`, `password_salt`, `password_type`, `role_id`, `org_cod`, `status`, `description`) VALUES 
(1, null, 1,      "p6m0adm",     1, "uBnKfXylWRdUFqVM424ERH.tISbfJbWq", "J3Apb1ZhNgtuBx4ifhg9F0MBVhI3bH9ELjJRQg==", "BCR3", null, "P6M", 1, "平台管理员"),
(2, null, 10007,  "it-10007",    1, "c557193f596ccf70b8cbc5ca306557b3", "uoqacs2t699ybv8tc42hz8z1shny6ups",         "MD5",  null, null,  1, "平台测试"),
(3, 1   , 1,      "13377777777", 2, null,     null, null, null, null,        1, "平台手机账户"),
(4, 1   , 1,      "f3@fmes.com", 3, null,     null, null, null, null,        1, "平台邮箱账户"),
(5, null, 100001, "it-100001",   1, "654321", null, null, null, null,        1, "赢迪管理员"),
(6, null, 100002, "it-100002",   1, "123456", null, null, null, "ORGCM3558", 1, "赢迪-测试1"),
(7, null, 100003, "it-100003",   1, "123456", null, null, null, "ORGCM3558", 1, "赢迪-测试1");

-- 平台服务
INSERT INTO `zgo_app_service`(`id`, `name`, `code`) VALUES 
(1, "认证授权", 'ka'),
(2, "灵活导购", 'lhdg'),
(3, "通路数据", 'tlsj');

-- 平台对应关系
INSERT INTO `zgo_app_service_audience`(`svc_id`, `audience`, `resource`) VALUES 
(1, null, '/api/ka/'),
(2, null, '/api/lhdg/'),
(3, null, '/api/tlsj/');

-- 授权
INSERT INTO `zgo_app_service_org`(`svc_id`, `org_cod`, `expired`, `status`, `description`) VALUES 
(1, "ORGCM3558", null, 1, ""),
(2, "ORGCM3558", null, 1, "");

-- domain 如果同一用户有多个角色,系统可以通过访问的域名自动分配角色
-- kid编码(24)： <时间编码8位><机器码4位><随机数12位>
INSERT INTO `zgo_role`(`id`, `kid`, `name`, `org_cod`, `org_adm`, `status`, `description`) VALUES
(1,  "202101121111000000000012", "Basis",           null,           0, 1, "系统基本访问权限"),
(2,  "202101121111000000000003", "Admin",          "ORGCM3558",     1, 1, "赢迪管理权限"),
(3,  "202101121111000000000004", "Organizer",      "ORGCM3558",     0, 1, "赢迪组织者权限"),
(4,  "202101121111000000000006", "Operator",       "ORGCM3558",     0, 1, "赢迪经营者权限"),
(5,  "202101121111000000000005", "OperatorOrg",    "ORGCM3558",     0, 1, "赢迪经营者(机构)权限"),
(6,  "202101121111000000000007", "OperatorField",  "ORGCM3558",     0, 1, "赢迪经营者(场)权限"),
(7,  "202101121111000000000008", "OperatorComm",   "ORGCM3558",     0, 1, "赢迪经营者(沟通)权限"),
(8,  "202101121111000000000009", "OperatorPoint",  "ORGCM3558",     0, 1, "赢迪经营者(支点)权限"),
(9,  "202101121111000000000010", "Executer",       "ORGCM3558",     0, 1, "赢迪执行者权限"),
(10, "202101121111000000000011", "Reseller",       "ORGCM3558",     0, 1, "赢迪经销商权限");

-- 角色继承
INSERT INTO `zgo_role_role`(`pid`, `cid`, `org_cod`) VALUES 
( 5,  6, 'ORGCM3558'),
( 5,  7, 'ORGCM3558'),
( 5,  8, 'ORGCM3558'),
( 5,  9, 'ORGCM3558');

-- 用户角色分配
INSERT INTO `zgo_user_role`(`user_id`, `role_id`, `org_cod`) VALUES 
(100002, 4, 'ORGCM3558'),
(100002, 5, 'ORGCM3558'),
(100003, 6, 'ORGCM3558');

-- 接口
INSERT INTO `zgo_policy_service_action`(`id`, `svc_id`, `name`, `resource`, `status`, `description`) VALUES 
(1, 1, "KaBasis",     "GET /*;/api/ka/v1/user/*",                                  1, "平台基本权限"),
(2, 1, "KaAdmin",     "/api/ka/*",                                                 1, "平台管理权限"),
(3, 2, "LhdgBasis",   "GET /*",                                                    1, "灵活导购基本权限"),
(4, 2, "LhdgAdmin",   "/api/lhdg/v1/*",                                            1, "灵活导购管理权限");

-- 策略
INSERT INTO `zgo_policy`(`id`, `org_cod`, `name`, `status`, `description`) VALUES 
(1, null, "LhdgAdmin", 1, "灵活导购管理权限策略"),
(2, null, "LhdgBasis", 1, "灵活导购基本权限策略"),
(3, "ORGCM3558", "Operator", 1, "赢迪经营者权限");

-- 策略声明
INSERT INTO `zgo_policy_statement`(`plcy_id`, `effect`, `action`, `resource`, `condition`) VALUES 
(1, 1, "lhdg:*",     null, null),
(2, 1, "lhdg:Basis", null, null),
(3, 1, "ka:KABasis;lhdg:Lhdg*",                                     null, null),
(3, 0, null, "fmes::lhdg:Organizer::/api/lhdg/v1/*", '{access_time":{"times": ["2020-12-12 00:00:00", "2021-02-02 00:00:00"]}}"');


-- 角色策略
INSERT INTO `zgo_role_policy`(`role_id`, `plcy_id`, `org_cod`) VALUES 
( 4,  3, 'ORGCM3558');

-- 角色策略
INSERT INTO `zgo_user_policy`(`user_id`, `plcy_id`, `org_cod`) VALUES 
(100003, 3, 'ORGCM3558');


-- 账户
INSERT INTO `zgo_account`(`id`, `pid`, `user_id`, `account`, `account_type`, `password`, `password_salt`, `password_type`, `role_id`, `org_cod`, `status`, `description`) VALUES
(8, null, 100003, "zgo1",   1, "123456", null, null, null, "ORGCM3558", 1, "ZGO-测试1")
(9, null, 100003, "plus1",   1, "123456", null, null, null, "ORGCM3558", 1, "PLUS-测试1")；



-- 编码说明
-- 用户(48)： u<助记符3位><时间编码8位><ID编码8位><机器码4位><随机码24位>
-- 租户(32)： o<助记符3位><时间编码8位><ID编码8位><机器码4位><随机码8位>
-- 应用(24)： a<助记符3位><时间编码8位><机器码4位><随机码8位>
INSERT INTO `zgo_user`(`id`, `type`, `kid`, `name`, `status`) VALUES
(100101, "org", "o00020210114xx1000011111000000a1",                 "测试",    1),
(100102, "usr", "u00020210114xx1000021111x000000000x000000000x002", "孟凡宇",  1),
(100103, "usr", "u00020210114xx1000031111x000000000x000000000x003", "曲海焦",  1),
(100104, "usr", "u00020210114xx1000031111x000000000x000000000x004", "孙媛媛",  1),
(100105, "usr", "u00020210114xx1000031111x000000000x000000000x005", "罗双慧",  1),
(100106, "usr", "u00020210114xx1000031111x000000000x000000000x006", "艾准",    1),
(100107, "usr", "u00020210114xx1000031111x000000000x000000000x007", "卢洪琦",  1),
(100108, "usr", "u00020210114xx1000031111x000000000x000000000x008", "庄园",    1),
(100109, "usr", "u00020210114xx1000031111x000000000x000000000x009", "洪波",    1),
(100110, "usr", "u00020210114xx1000031111x000000000x000000000x010", "邝艳秋",  1),
(100111, "usr", "u00020210114xx1000031111x000000000x000000000x011", "童静",    1),
(100112, "usr", "u00020210114xx1000031111x000000000x000000000x012", "陈哲",    1),
(100113, "usr", "u00020210114xx1000031111x000000000x000000000x013", "李丹丹",  1),
(100114, "usr", "u00020210114xx1000031111x000000000x000000000x014", "李根",    1),
(100115, "usr", "u00020210114xx1000031111x000000000x000000000x015", "艾磊",    1);

-- 租户用户
-- 租户编码(40)： u<助记符3位><租户ID编码8位><时间编码8位><用户ID编码8位><机器码4位><随机码8位>
INSERT INTO `zgo_organization_user`(`user_id`, `org_cod`, `ouid`, `name`) VALUES
(100002, "ORGCM3558", "u000xx10000120210112xx100002111100000001", "赢迪测试2"),
(100003, "ORGCM3558", "u000xx10000120210112xx100003111100000001", "赢迪测试3");

1.0编码	经营者	扮演的人的名字
NSY20102515	孟凡宇	李丹
PT20190124	曲海焦	徐薇
HSH20180249	孙媛媛	游旋
HSH20160180	罗双慧	罗双慧
HSH20070123	艾准	刘静1
HSH20151304	卢洪琦	卢洪琦
OT202000002	庄园	庄园
PTB0420151382	洪波	吴小英
HSH20160589	秋	邝艳秋
SPOC20180044	童静	王明
PTB0420140083	陈哲	程慧
PTB0220161379	李丹丹	徐浩
PTB0720160124	李根	杨小慧
PTB0920151531	艾磊	李瑾
