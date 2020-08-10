package service

import (
	"errors"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/google/wire"
	zgocasbin "github.com/suisrc/zgo/modules/casbin"
)

// CasbinAdapterSet 注入casbin
var CasbinAdapterSet = wire.NewSet(
	zgocasbin.NewCasbinEnforcer,
	wire.Struct(new(CasbinAdapter), "GPA"),
	// NewCasbinAdapter,
	wire.Bind(new(persist.Adapter), new(CasbinAdapter)),
)

// CasbinAdapter 账户管理
type CasbinAdapter struct {
	GPA              // 数据库
	VerPolicy string // adapter版本,防止重复更新
}

// ================================================ 分割线

var _ zgocasbin.PolicyVer = (*CasbinAdapter)(nil)

// PolicyVer ver
func (a *CasbinAdapter) PolicyVer() string {
	return a.VerPolicy
}

// PolicySet set
func (a *CasbinAdapter) PolicySet(ver string) error {
	a.VerPolicy = ver
	return nil
}

// ================================================ 分割线

var _ persist.Adapter = (*CasbinAdapter)(nil)

// LoadPolicy loads policy from database.
func (a *CasbinAdapter) LoadPolicy(model model.Model) error {
	// resouces

	// role

	// role-role

	persist.LoadPolicyLine()
	return nil
}

// SavePolicy saves policy to database.
func (a *CasbinAdapter) SavePolicy(model model.Model) error {
	// 该方法只要通知系统接受更新
	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
