// +build document

package casbinredis

//import (
//	"github.com/casbin/casbin/v2"
//	"github.com/casbin/casbin/v2/persist"
//	zgocasbin "github.com/suisrc/zgo/casbin"
//	"github.com/suisrc/zgo/
//
//	rediswatcher "github.com/billcobbler/casbin-redis-watcher/v2"
//)

// NewCasbinWatcher 构建Casbin Watcher
// 启动观察者模式,可以通过enforcer.SavePolicy, 通知集群更新policy内容
// func NewCasbinWatcher(pv zgocasbin.PolicyVer, enforcer *casbin.SyncedEnforcer) (persist.Watcher, error) {
// 	enforcer.EnableAutoSave(false)
// 	// ctx, cancel := context.WithCancel(context.Background())
// 	// defer cancel()
// 	watcher, err := rediswatcher.NewWatcher("redis-svc:6379")
// 	if err != nil {
// 		return nil, err
// 	}
// 	enforcer.SetWatcher(watcher)
// 	watcher.SetUpdateCallback(func(ver string) {
// 		if pv.PolicyVer() != ver { // 控制执行的版本
// 			enforcer.LoadPolicy()
// 			pv.PolicySet(ver)
// 			logger.Infof(nil, "reloading casbin: version => %s", ver)
// 		}
// 	})
// 	return watcher, nil
// }
