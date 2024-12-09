package red_lock

import (
	"context"
	"gin-ddd-example/example/redis_example"
	"gin-ddd-example/internal/app/repo"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/stretchr/testify/suite"
)

type RedLockTestSuite struct {
	redis_example.RedisTestSuite
	productStockRepo *repo.ProductStockRepoImpl
}

func TestRedLockTestSuite(t *testing.T) {
	suite.Run(t, new(RedLockTestSuite))
}

func (suite *RedLockTestSuite) SetupTest() {
	suite.RedisTestSuite.SetupTest()
	suite.productStockRepo = repo.NewProductStockRepo(suite.Db)
}

func (suite *RedLockTestSuite) TearDownTest() {
	suite.Rdb.Close()
}

// 基本使用
func (suite *RedLockTestSuite) TestRedLock_Scenario1() {
	pool := goredis.NewPool(suite.Rdb)
	// 创建一个 redsync 实例
	rs := redsync.New(pool)

	// 创建一个互斥锁
	mutex := rs.NewMutex("test-redlock",
		redsync.WithExpiry(time.Second*5),
		redsync.WithTries(3),
	)

	// 获取锁
	if err := mutex.Lock(); err != nil {
		suite.T().Fatalf("获取锁失败: %v", err)
	}
	// 释放锁
	defer mutex.Unlock()
	// 处理业务逻辑
	suite.T().Log("获取锁成功")
}

// 带超时的锁操作
func (suite *RedLockTestSuite) TestRedisLock_Scenario2() {
	pool := goredis.NewPool(suite.Rdb)
	rs := redsync.New(pool)
	mutex := rs.NewMutex("my-key",
		// 自定义选项
		redsync.WithExpiry(time.Second*8),
		redsync.WithTries(3),
		redsync.WithRetryDelay(time.Millisecond*100),
		redsync.WithDriftFactor(0.01),
		redsync.WithTimeoutFactor(0.05),
	)
	// 使用上下文控制超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := mutex.LockContext(ctx); err != nil {
		suite.T().Fatalf("获取锁失败: %v", err)
		return
	}
	defer mutex.UnlockContext(ctx)
}
