package schema

// api调用service时候,中间传递的结构体
const (
	TablePrefix = "zgo_"
	WhereIS     = false
)

// AccountType 账户类型
type AccountType int

const (
	// ATNone 无
	ATNone AccountType = iota // value -> 0
	// ATName 名称
	ATName
	// ATMobile 手机
	ATMobile
	// ATEmail 邮箱
	ATEmail
	// ATOpenid openid
	ATOpenid
	// ATUnionid unionid
	ATUnionid
	// ATToken token
	ATToken
)
