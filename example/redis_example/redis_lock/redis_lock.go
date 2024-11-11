package redis_lock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	rdb        *redis.Client
	key        string
	value      string        // 锁的持有者标识
	expiration time.Duration // 锁的过期时间
	stopChan   chan struct{} // 用于停止续租
}

func NewRedisLock(rdb *redis.Client, key string, value string, expiration time.Duration) *RedisLock {
	return &RedisLock{
		rdb:        rdb,
		key:        key,
		value:      value,
		expiration: expiration,
	}
}

// 尝试获取锁
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
	success, err := l.rdb.SetNX(ctx, l.key, l.value, l.expiration).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

// 释放锁，只有锁的持有者才能释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	result, err := l.rdb.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return errors.New("锁不属于当前持有者")
	}
	return nil
}

// 检查是否仍然持有锁
func (l *RedisLock) CheckLockOwner(ctx context.Context) bool {
	val, err := l.rdb.Get(ctx, l.key).Result()
	if err != nil {
		return false
	}
	return val == l.value
}

// 开启锁的续租
func (l *RedisLock) EnableAutoRenew(ctx context.Context) {
	l.stopChan = make(chan struct{})
	go func() {
		ticker := time.NewTicker(l.expiration / 3) // 在过期时间的1/3处续租
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 检查并续租
				success, err := l.Renew(ctx)
				if err != nil || !success {
					return
				}
			case <-l.stopChan:
				return
			}
		}
	}()
}

// 续租
func (l *RedisLock) Renew(ctx context.Context) (bool, error) {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
	`
	result, err := l.rdb.Eval(ctx, script,
		[]string{l.key}, l.value, l.expiration.Microseconds()).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

func (l *RedisLock) StopAutoRenew() {
	if l.stopChan != nil {
		close(l.stopChan)
	}
}
