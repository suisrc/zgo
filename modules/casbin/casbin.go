package casbin

import (
	"errors"
	"strings"
	"time"

	"github.com/suisrc/zgo/modules/logger"

	"github.com/suisrc/zgo/modules/config"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
)

// PolicyVer 政策版本信息
type PolicyVer interface {
	PolicyVer() string
	PolicySet(string) error
}

// NewCasbinEnforcer 初始化casbin enforcer
func NewCasbinEnforcer(adapter persist.Adapter) (*casbin.SyncedEnforcer, func(), error) {
	c := config.C.Casbin
	if c.Model == "" {
		// return new(casbin.SyncedEnforcer), func() {}, nil
		return nil, nil, errors.New("casbin model no config")
	}

	enforcer, err := casbin.NewSyncedEnforcer(c.Model)
	if err != nil {
		return nil, nil, err
	}
	logger.Infof(nil, "loading casbin model[%s]", c.Model)
	enforcer.EnableLog(c.Debug)

	err = enforcer.InitWithModelAndAdapter(enforcer.GetModel(), adapter)
	if err != nil {
		return nil, nil, err
	}
	enforcer.EnableEnforce(c.Enable)

	cleanFunc := func() {}
	if c.AutoLoad {
		enforcer.StartAutoLoadPolicy(time.Duration(c.AutoLoadInternal) * time.Second)
		cleanFunc = func() {
			enforcer.StopAutoLoadPolicy()
		}
	}

	// 注册方法
	enforcer.AddFunction("fdom", DomainMatchFunc)
	enforcer.AddFunction("fact", ActionMatchFunc)
	enforcer.AddFunction("fdoma", DomainMatchAudienceFunc)

	return enforcer, cleanFunc, nil
}

//====================================
// func
//====================================

// DomainMatch domain
func DomainMatch(key1 string, key2 string) bool {
	if key2[:1] == "." {
		return strings.HasSuffix(key1, key2)
	}
	i := strings.Index(key2, "*") + 1
	if i == 0 {
		return key1 == key2
	}
	l := len(key2)
	if i == l {
		return true
	}
	if li := len(key1) - (l - i); li > 0 {
		// 截取key1可用部分
		return key1[li:] == key2[i:]
	}
	return key1 == key2[i:]
}

// DomainMatchFunc domain
func DomainMatchFunc(args ...interface{}) (interface{}, error) {
	domain1 := args[0].(string)
	domain2 := args[1].(string)
	if domain2 == "" || domain2 == "*" {
		return true, nil
	}
	return DomainMatch(domain1, domain2), nil
}

// DomainMatchAudienceFunc domain
func DomainMatchAudienceFunc(args ...interface{}) (interface{}, error) {
	domain1 := args[0].(string)
	domain2 := args[1].(string)
	audience := args[2].(string)
	if domain2 == "" || domain2 == "*" {
		return true, nil
	}
	if domain2 == "jwt" {
		return DomainMatch(domain1, audience), nil
	}
	return DomainMatch(domain1, domain2), nil
}

// ActionMatchFunc action
func ActionMatchFunc(args ...interface{}) (interface{}, error) {
	action := args[0].(string)
	actions := args[1].(string)
	if actions == "" || actions == "*" {
		return true, nil
	}
	return strings.Contains(actions, "("+action+")"), nil
}
