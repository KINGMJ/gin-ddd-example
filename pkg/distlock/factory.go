package distlock

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// LockerFactory 锁工厂实现
type LockerFactory struct {
	redisClient *redis.Client
	etcdClient  *clientv3.Client
}

// NewLockerFactory 创建锁工厂
func NewLockerFactory(redisClient *redis.Client, etcdClient *clientv3.Client) *LockerFactory {
	return &LockerFactory{
		redisClient: redisClient,
		etcdClient:  etcdClient,
	}
}

func (f *LockerFactory) Create(opts ...OptionFunc) (Locker, error) {
	options := DefaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	switch options.LockerType {
	case LockerTypeRedSync:
		return NewRedSyncLocker(f.redisClient, options), nil
	case LockerTypeEtcd:
		return NewEtcdLocker(f.etcdClient, options), nil
	default:
		return nil, fmt.Errorf("unsupported locker type: %s", options.LockerType)
	}
}
