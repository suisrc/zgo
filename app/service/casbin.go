package service

import (
	"errors"

	"github.com/suisrc/zgo/middleware"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/google/wire"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/schema"
	zgocasbin "github.com/suisrc/zgo/modules/casbin"
	"github.com/suisrc/zgo/modules/config"
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
	nosignin, norole, nouser := false, false, true
	// resouces
	resource0 := schema.CasbinGpaResource{}
	resources := []schema.CasbinGpaResource{}
	if err := resource0.QueryAll(a.Sqlx, &resources); err != nil {
		logger.Infof(nil, "loading casbin: none -> %s", logger.ErrorWW(err))
		return nil
	}
	for _, r := range resources {
		if !r.Resource.Valid {
			continue
		}
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
		if !nosignin && r.Resource.String == middleware.CasbinNoSignin {
			nosignin = true
		} else if !norole && r.Resource.String == middleware.CasbinNoRole {
			norole = true
		}
	}
	// user
	user0 := schema.CasbinGpaResourceUser{}
	users := []schema.CasbinGpaResourceUser{}
	if err := user0.QueryAll(a.Sqlx, &users); err != nil && !sqlxc.IsNotFound(err) {
		logger.Infof(nil, "loading casbin: user -> %s", logger.ErrorWW(err))
	}
	for _, r := range users {
		if !r.User.Valid || !r.Resource.Valid {
			continue
		}
		line := "g2"
		line += "," + middleware.CasbinUserPrefix + r.User.String
		line += "," + r.Resource.String
		persist.LoadPolicyLine(line, model)
		logger.Infof(nil, "loading casbin: %s", line)
		if nouser {
			nouser = false
		}
	}
	// config
	config.C.Casbin.NoSignin = nosignin // 覆盖性修改默认配置
	config.C.Casbin.NoRole = norole     // 覆盖性修改默认配置
	config.C.Casbin.NoUser = nouser     // 覆盖性修改默认配置
	logger.Infof(nil, "loading casbin: nosignin: %t, norole: %t, nouser: %t", nosignin, norole, nouser)
	// role
	role0 := schema.CasbinGpaResourceRole{}
	roles := []schema.CasbinGpaResourceRole{}
	if err := role0.QueryAll(a.Sqlx, &roles); err != nil {
		if !sqlxc.IsNotFound(err) {
			logger.Infof(nil, "loading casbin: role -> %s", logger.ErrorWW(err))
		}
		return nil
	}
	for _, r := range roles {
		if !r.Role.Valid || !r.Resource.Valid {
			continue
		}
		line := "g"
		line += "," + middleware.CasbinRolePrefix + r.Role.String
		line += "," + r.Resource.String
		persist.LoadPolicyLine(line, model)
		logger.Infof(nil, "loading casbin: %s", line)
	}
	// role-role
	rolerole0 := schema.CasbinGpaRoleRole{}
	roleroles := []schema.CasbinGpaRoleRole{}
	if err := rolerole0.QueryAll(a.Sqlx, &roleroles); err != nil {
		if !sqlxc.IsNotFound(err) {
			logger.Infof(nil, "loading casbin: role-role -> %s", logger.ErrorWW(err))
		}
		return nil
	}
	for _, r := range roleroles {
		if !r.Owner.Valid || !r.Child.Valid {
			continue
		}
		line := "g"
		line += "," + middleware.CasbinRolePrefix + r.Owner.String
		line += "," + middleware.CasbinRolePrefix + r.Child.String
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
