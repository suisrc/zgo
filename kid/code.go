package kid

import (
	"strings"
)

/*
用户 user_id=1，为超级管理员，所有的权限认证将被跳过, 系统初始化完成后， 最好将该角色权限禁用
*/

// NewRoleKID ... 角色(24) r<助记符3位><时间编码8位><机器码4位><随机码8位>
func NewRoleKID() string {
	var builder strings.Builder
	builder.WriteString("rPv2")
	builder.WriteString(NewNowCode(8, false))
	builder.WriteString(GetMustMachineCode())
	builder.WriteString(NewSequenceCode(8, true))
	return builder.String()
}

// NewUsrKID ... 用户(36) u<助记符3位><时间编码8位><ID编码8位><机器码4位><随机码12位>
func NewUsrKID(id int64) string {
	var builder strings.Builder
	builder.WriteString("uPv2")
	builder.WriteString(NewNowCode(8, false))
	builder.WriteString(NewIdxCode(8, id, true))
	builder.WriteString(GetMustMachineCode())
	builder.WriteString(NewSequenceCode(12, true))
	return builder.String()
}

// NewOrgKID ... 租户(32) t<助记符3位><时间编码8位><ID编码8位><机器码4位><随机码8位>
func NewOrgKID(id int64) string {
	var builder strings.Builder
	builder.WriteString("tPv2")
	builder.WriteString(NewNowCode(8, false))
	builder.WriteString(NewIdxCode(8, id, true))
	builder.WriteString(NewSequenceCode(8, true))
	return builder.String()
}

// NewAppKID ... 应用(24) a<助记符3位><时间编码8位><机器码4位><随机码8位>
func NewAppKID() string {
	var builder strings.Builder
	builder.WriteString("aPv2")
	builder.WriteString(NewNowCode(8, false))
	builder.WriteString(GetMustMachineCode())
	builder.WriteString(NewSequenceCode(8, true))
	return builder.String()
}
