// +build document

package casbinmem

//import (
//	"context"
//
//	"github.com/casbin/casbin/v2"
//	"github.com/casbin/casbin/v2/persist"
//	cloudwatcher "github.com/rusenask/casbin-go-cloud-watcher"
//	zgocasbin "github.com/suisrc/zgo/casbin"
//
//	// Enable in-memory driver
//	_ "github.com/rusenask/casbin-go-cloud-watcher/drivers/mempubsub"
//)

// NewCasbinWatcher 构建Casbin Watcher
// 启动观察者模式,可以通过enforcer.SavePolicy, 通知集群更新policy内容
// func NewCasbinWatcher(pv zgocasbin.PolicyVer, enforcer *casbin.SyncedEnforcer) (persist.Watcher, func(), error) {
// 	enforcer.EnableAutoSave(false)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	// defer cancel()
//
// 	watcher, err := cloudwatcher.New(ctx, "mem://topicA")
// 	if err != nil {
// 		// 无法启动监听同步
// 		cancel()
// 		return nil, nil, err
// 	}
// 	enforcer.SetWatcher(watcher)
// 	watcher.SetUpdateCallback(func(ver string) {
// 		if pv.PolicyVer() != ver { // 控制执行的版本
// 			enforcer.LoadPolicy()
// 			pv.PolicySet(ver)
// 			logger.Infof(nil, "reloading casbin: version => %s", ver)
// 		}
// 	})
// 	return watcher, cancel, nil
// }
