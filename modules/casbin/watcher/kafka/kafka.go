package casbinkafka

//import (
//	"context"
//	"os"
//
//	"github.com/casbin/casbin/v2"
//	"github.com/casbin/casbin/v2/persist"
//	cloudwatcher "github.com/rusenask/casbin-go-cloud-watcher"
//	zgocasbin "github.com/suisrc/zgo/modules/casbin"
//	"github.com/suisrc/zgo/modules/logger"
//
//	// Enable in-memory driver
//	_ "github.com/rusenask/casbin-go-cloud-watcher/drivers/kafkapubsub"
//)

// NewCasbinWatcher 构建Casbin Watcher
// 启动观察者模式,可以通过enforcer.SavePolicy, 通知集群更新policy内容
// func NewCasbinWatcher(pv zgocasbin.PolicyVer, enforcer *casbin.SyncedEnforcer) (persist.Watcher, func(), error) {
// 	enforcer.EnableAutoSave(false)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	// defer cancel()
//
// 	// Watcher can publish to a Kafka cluster.
// 	// A Kafka URL only includes the topic name.
// 	// The brokers in the Kafka cluster are discovered from the KAFKA_BROKERS environment variable
// 	// (which is a comma-delimited list of hosts, something like 1.2.3.4:9092,5.6.7.8:9092).
//
// 	// The set of brokers must be in an environment variable KAFKA_BROKERS.
// 	os.Setenv("KAFKA_BROKERS", "kafka-svc:9092")
// 	watcher, err := cloudwatcher.New(ctx, "kafka://topicA")
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
