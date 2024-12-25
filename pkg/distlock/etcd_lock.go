package distlock

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// EtcdLocker Etcd 分布式锁实现
type EtcdLocker struct {
	client  *clientv3.Client
	options Options
}

type etcdLock struct {
	session *concurrency.Session
	mutex   *concurrency.Mutex
}

func NewEtcdLocker(client *clientv3.Client, options Options) Locker {
	return &EtcdLocker{
		client:  client,
		options: options,
	}
}

func (l *EtcdLocker) Lock(ctx context.Context, key string) (Lock, error) {
	// 创建 session
	expiry := int(l.options.Expiry.Seconds())
	session, err := concurrency.NewSession(l.client, concurrency.WithTTL(expiry))
	if err != nil {
		return nil, fmt.Errorf("failed to new etcd session: %w", err)
	}
	// 创建 mutex
	key = l.options.KeyPrefix + key
	mutex := concurrency.NewMutex(session, key)

	// 尝试获取锁
	if err = mutex.Lock(ctx); err != nil {
		closeErr := session.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to close etcd session: %w", closeErr)
		}
		return nil, fmt.Errorf("failed to acquire lock for key %s: %w", key, err)
	}
	return &etcdLock{
		session: session,
		mutex:   mutex,
	}, nil
}

func (l *etcdLock) Unlock(ctx context.Context) (err error) {
	defer func() {
		closeErr := l.session.Close()
		if closeErr != nil {
			err = fmt.Errorf("failed to close etcd session: %w", closeErr)
		}
	}()
	if err = l.mutex.Unlock(context.Background()); err != nil {
		err = fmt.Errorf("failed to unlock etcd lock: %w", err)
	}
	return
}
