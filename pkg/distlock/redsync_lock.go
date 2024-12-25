package distlock

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// ErrLockNotHeld 表示尝试释放一个未持有的锁
var ErrLockNotHeld = errors.New("lock not held")

type RedSyncLocker struct {
	rs      *redsync.Redsync
	options Options
}

type redSyncLock struct {
	mutex *redsync.Mutex
}

func NewRedSyncLocker(client *redis.Client, options Options) Locker {
	pool := goredis.NewPool(client)
	// 创建一个 redsync 实例
	rs := redsync.New(pool)
	return &RedSyncLocker{
		rs:      rs,
		options: options,
	}
}

func (l *RedSyncLocker) Lock(ctx context.Context, key string) (Lock, error) {
	key = l.options.KeyPrefix + key
	mutex := l.rs.NewMutex(key,
		redsync.WithExpiry(l.options.Expiry),
		redsync.WithTries(l.options.RetryCount),
		redsync.WithRetryDelay(l.options.RetryDelay),
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to acquire lock for key %s: %w", key, err)
	}
	return &redSyncLock{mutex: mutex}, nil
}

func (l *redSyncLock) Unlock(ctx context.Context) error {
	ok, err := l.mutex.UnlockContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to unlock lock: %w", err)
	}
	if !ok {
		return ErrLockNotHeld
	}
	return nil
}
