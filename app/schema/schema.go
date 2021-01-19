package schema

// api调用service时候,中间传递的结构体
const (
	TablePrefix = "zgo_"
	WhereIS     = false
)

// AccountType 账户类型
type AccountType int

const (
	// AccountTypeNone 无
	AccountTypeNone AccountType = iota // value -> 0
	// AccountTypeName 名称
	AccountTypeName
	// AccountTypeMobile 手机
	AccountTypeMobile
	// AccountTypeEmail 邮箱
	AccountTypeEmail
	// AccountTypeOpenid openid
	AccountTypeOpenid
	// AccountTypeUnionid unionid
	AccountTypeUnionid
	// AccountTypeToken token
	AccountTypeToken
)

// StatusType 数据状态
type StatusType int

// 1:启用 0:禁用 2: 未激活 3: 注销
const (
	// StatusDisable 禁用
	StatusDisable StatusType = iota // value -> 0
	// StatusEnable 启用
	StatusEnable
	// StatusNoActivate 为激活
	StatusNoActivate
	// StatusRevoked 注销
	StatusRevoked
	// StatusDeleted 删除
	StatusDeleted
)
