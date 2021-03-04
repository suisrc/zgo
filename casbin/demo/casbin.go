// +build document

package demo

// PolicyVer 政策版本信息
// type PolicyVer interface {
// 	PolicyVer() string
// 	PolicySet(string) error
// }
//
// // NewCasbinEnforcer 初始化casbin enforcer
// func NewCasbinEnforcer(adapter persist.Adapter) (*casbin.SyncedEnforcer, func(), error) {
// 	c := config.C.Casbin
// 	if c.Model == "" {
// 		// return new(casbin.SyncedEnforcer), func() {}, nil
// 		return nil, nil, errors.New("casbin model no config")
// 	}
//
// 	enforcer, err := casbin.NewSyncedEnforcer(c.Model, adapter)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	logger.Infof(nil, "loading casbin model[%s]", c.Model)
// 	enforcer.EnableLog(c.Debug)
// 	enforcer.EnableEnforce(c.Enable)
//
// 	cleanFunc := func() {}
// 	if c.AutoLoad {
// 		enforcer.StartAutoLoadPolicy(time.Duration(c.AutoLoadInternal) * time.Second)
// 		cleanFunc = func() {
// 			enforcer.StopAutoLoadPolicy()
// 		}
// 	}
//
// 	// 注册方法
// 	enforcer.AddFunction("domainMatch", DomainMatchFunc)
// 	enforcer.AddFunction("methodMatch", MethodMatchFunc)
// 	enforcer.AddFunction("audienceMatch", AudienceMatchFunc)
//
// 	return enforcer, cleanFunc, nil
// }
