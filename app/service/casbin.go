package service

import (
	"errors"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/google/wire"
	"github.com/suisrc/zgo/app/schema"
	zgocasbin "github.com/suisrc/zgo/modules/casbin"
	"github.com/suisrc/zgo/modules/logger"
)

// CasbinAdapterSet 注入casbin
var CasbinAdapterSet = wire.NewSet(
	zgocasbin.NewCasbinEnforcer,
	wire.Struct(new(CasbinAdapter), "GPA"),
	wire.Bind(new(persist.Adapter), new(CasbinAdapter)),
	// NewCasbinAdapter,
)

// CasbinAdapter 账户管理
type CasbinAdapter struct {
	GPA              // 数据库
	VerPolicy string // adapter版本,防止重复更新
}

// ================================================ 分割线

var _ zgocasbin.PolicyVer = (*CasbinAdapter)(nil)

// PolicyVer ver
func (a CasbinAdapter) PolicyVer() string {
	return a.VerPolicy
}

// PolicySet set
func (a CasbinAdapter) PolicySet(ver string) error {
	a.VerPolicy = ver
	return nil
}

// ================================================ 分割线

var _ persist.Adapter = (*CasbinAdapter)(nil)

// LoadPolicy loads policy from database.
func (a CasbinAdapter) LoadPolicy(model model.Model) error {
	// resouces
	resource0 := schema.CasbinGpaResource{}
	resources := []schema.CasbinGpaResource{}
	err := a.GPA.Sqlx.Select(&resources, resource0.SQLByALL())
	if err != nil {
		logger.Infof(nil, "loading casbin: none -> %s", err.Error())
		return nil
	}
	for _, r := range resources {
		line := "p"
		line += "," + r.Resource.String
		line += "," + r.Domain.String
		line += "," + r.Path.String
		line += "," + r.Netmask.String
		line += "," + r.Methods.String
		if r.Allow.Bool {
			line += ",allow"
		} else {
			line += ",deny"
		}
		persist.LoadPolicyLine(line, model)
		logger.Infof(nil, "loading casbin: %s", line)
	}
	// role
	role0 := schema.CasbinGpaResourceRole{}
	roles := []schema.CasbinGpaResourceRole{}
	err = a.GPA.Sqlx.Select(&roles, role0.SQLByALL())
	if err != nil {
		return nil
	}
	for _, r := range roles {
		line := "g"
		line += "," + r.Role.String
		line += "," + r.Resource.String
		persist.LoadPolicyLine(line, model)
		logger.Infof(nil, "loading casbin: %s", line)
	}
	// role-role
	rolerole0 := schema.CasbinGpaRoleRole{}
	roleroles := []schema.CasbinGpaRoleRole{}
	err = a.GPA.Sqlx.Select(&roleroles, rolerole0.SQLByALL())
	if err != nil {
		return nil
	}
	for _, r := range roleroles {
		line := "g"
		line += "," + r.Owner.String
		line += "," + r.Child.String
		persist.LoadPolicyLine(line, model)
		logger.Infof(nil, "loading casbin: %s", line)
	}

	return nil
}

// SavePolicy saves policy to database.
func (a CasbinAdapter) SavePolicy(model model.Model) error {
	// 该方法只要通知系统接受更新, 无处理任何内容
	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
