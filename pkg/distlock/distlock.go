package distlock

import (
	"context"
	"time"
)

// 分布式锁实现

// Locker 分布式锁接口
type Locker interface {
	// Lock 获取锁
	Lock(ctx context.Context, key string) (Lock, error)
}

// Lock 锁接口
type Lock interface {
	// Unlock 释放锁
	Unlock(ctx context.Context) error
}

// Factory 锁工厂接口
type Factory interface {
	// Create 创建锁实例
	Create(opts ...OptionFunc) (Locker, error)
}

// 锁的实现类型
type LockerType string

const (
	LockerTypeRedSync LockerType = "redsync"
	LockerTypeEtcd    LockerType = "etcd"
)

// Options 分布式锁配置选项
type Options struct {
	LockerType LockerType
	Expiry     time.Duration
	RetryCount int
	RetryDelay time.Duration
	KeyPrefix  string
}

// OptionFunc 函数类型，用于配置定义Options选项
// 这是一种设计模式：Functional Options Pattern
type OptionFunc func(*Options)

// DefaultOptions 默认配置
var DefaultOptions = Options{
	LockerType: LockerTypeRedSync,
	Expiry:     10 * time.Second,
	RetryCount: 3,
	RetryDelay: 100 * time.Millisecond,
	KeyPrefix:  "lock:",
}

// 设置LockerType
func WithType(t LockerType) OptionFunc {
	return func(o *Options) {
		o.LockerType = t
	}
}

func WithExpiry(expiry time.Duration) OptionFunc {
	return func(o *Options) {
		o.Expiry = expiry
	}
}

func WithRetryCount(count int) OptionFunc {
	return func(o *Options) {
		o.RetryCount = count
	}
}

func WithRetryDelay(delay time.Duration) OptionFunc {
	return func(o *Options) {
		o.RetryDelay = delay
	}
}

func WithKeyPrefix(prefix string) OptionFunc {
	return func(o *Options) {
		o.KeyPrefix = prefix
	}
}
