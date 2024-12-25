package distlock

import (
	"gin-ddd-example/pkg/config"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// 分布式锁实例化
func NewDistLock(c config.Config, rdb *redis.Client, etcdClient *clientv3.Client) Locker {
	factory := NewLockerFactory(rdb, etcdClient)
	locker, err := factory.Create(
		WithType(LockerType(c.DistLock.Type)),
		WithRetryCount(c.DistLock.RetryCount),
		WithRetryDelay(time.Duration(c.DistLock.RetryDelay)*time.Millisecond),
	)
	if err != nil {
		panic(err)
	}
	return locker
}
