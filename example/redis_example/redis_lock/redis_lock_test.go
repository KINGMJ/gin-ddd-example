package redis_lock

import (
	"errors"
	"fmt"
	"gin-ddd-example/example/redis_example"
	"gin-ddd-example/internal/app/repo"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type RedisLockTestSuite struct {
	redis_example.RedisTestSuite
	productStockRepo *repo.ProductStockRepoImpl
}

func TestRedisLockTestSuite(t *testing.T) {
	suite.Run(t, new(RedisLockTestSuite))
}

func (suite *RedisLockTestSuite) SetupTest() {
	suite.RedisTestSuite.SetupTest()
	suite.productStockRepo = repo.NewProductStockRepo(suite.Db)
}

// 模拟三个用户并发下单
func (suite *RedisLockTestSuite) TestRedisLock_Scenario1() {
	// 模拟用户下单购买流程
	var productID int64 = 123599
	usersCount := 3 // 三个用户并发下单

	var wg sync.WaitGroup
	wg.Add(usersCount)

	for i := 0; i < usersCount; i++ {
		go func(userID int) {
			defer wg.Done()
			// 获取 redis 锁
			lock := NewRedisLock(suite.Rdb,
				fmt.Sprintf("product_lock_%d", productID),
				fmt.Sprintf("%d", userID),
				time.Second*5,
			)
			ok, err := lock.TryLock(suite.Ctx)
			if err != nil || !ok {
				suite.T().Logf("用户%d获取锁失败", userID)
				return
			}
			defer lock.Unlock(suite.Ctx) // 释放锁

			// 模拟下单流程
			suite.T().Logf("用户%d开始下单...", userID)
			err = suite.Db.WithContext(suite.Ctx).Transaction(func(tx *gorm.DB) error {
				err := suite.productStockRepo.DeductStock(tx, productID, 1)
				if err != nil {
					return err
				}
				suite.T().Logf("用户%d下单成功，扣减%d件库存", userID, 1)
				return nil
			})
			if err != nil {
				suite.T().Logf("用户%d下单失败，%v", userID, err)
			}
		}(i + 1)
	}

	// 等待所有用户下单完成
	wg.Wait()

	// 最终验证
	finalStock, err := suite.productStockRepo.FindByProductId(productID)
	suite.NoError(err)
	suite.T().Logf("最终库存：%d", finalStock.Count)
}

// 模拟三个用户下单，用户1获取锁后崩溃，用户2立即尝试获取锁失败，用户3等待锁过期后获取锁
func (suite *RedisLockTestSuite) TestRedisLock_Scenario2() {
	var productID int64 = 123599

	var wg sync.WaitGroup
	wg.Add(3)

	// 模拟用户1： 获取锁然后崩溃
	go func() {
		defer wg.Done()
		lock1 := NewRedisLock(suite.Rdb,
			fmt.Sprintf("product_lock_%d", productID),
			fmt.Sprintf("%d", 1),
			time.Second*5,
		)
		ok, err := lock1.TryLock(suite.Ctx)
		suite.NoError(err)
		suite.True(ok, "客户端1应该能够获取到锁")
		// 模拟客户端1崩溃，未释放锁...
		suite.T().Log("客户端1崩溃，未释放锁...")
	}()

	time.Sleep(time.Millisecond * 100)
	// 模拟客户端2：立即尝试获取锁（应该失败）
	go func() {
		defer wg.Done()
		lock2 := NewRedisLock(suite.Rdb,
			fmt.Sprintf("product_lock_%d", productID),
			fmt.Sprintf("%d", 2),
			time.Second*5,
		)
		ok, _ := lock2.TryLock(suite.Ctx)
		suite.False(ok, "客户端2应该无法获取到锁")
		suite.T().Log("客户端2尝试获取锁失败（锁未过期）")
	}()

	// 等待锁过期
	time.Sleep(5 * time.Second)
	go func() {
		defer wg.Done()
		lock3 := NewRedisLock(suite.Rdb,
			fmt.Sprintf("product_lock_%d", productID),
			fmt.Sprintf("%d", 3),
			time.Second*5,
		)
		ok, err := lock3.TryLock(suite.Ctx)
		suite.NoError(err)
		suite.True(ok, "锁已过期，客户端3应该能获取到锁")
		// 模拟下单流程
		suite.T().Logf("用户%d开始下单...", 3)
		err = suite.Db.WithContext(suite.Ctx).Transaction(func(tx *gorm.DB) error {
			err := suite.productStockRepo.DeductStock(tx, productID, 1)
			if err != nil {
				return err
			}
			suite.T().Logf("用户%d下单成功，扣减%d件库存", 3, 1)
			return nil
		})
		if err != nil {
			suite.T().Logf("用户%d下单失败，%v", 3, err)
		}
	}()

	wg.Wait()
}

// 模拟三个用户下单，用户1获取锁后执行耗时操作，用户2立即尝试获取锁失败，用户3等待锁过期后获取锁
// 需要加入锁的持有者检查
func (suite *RedisLockTestSuite) TestRedisLock_Scenario3() {
	var productID int64 = 123599
	lockTimeout := time.Second * 5

	var wg sync.WaitGroup
	wg.Add(2)

	// 模拟用户1：获取锁后执行耗时操作
	go func() {
		defer wg.Done()
		lock1 := NewRedisLock(suite.Rdb,
			fmt.Sprintf("product_lock_%d", productID),
			"client_1", // 使用唯一的客户端标识
			lockTimeout,
		)

		ok, err := lock1.TryLock(suite.Ctx)
		suite.NoError(err)
		suite.True(ok, "客户端1应该能够获取到锁")
		suite.T().Log("客户端1获取锁成功")

		// 模拟耗时操作
		time.Sleep(lockTimeout + time.Second)

		// 尝试执行业务操作
		err = suite.Db.WithContext(suite.Ctx).Transaction(func(tx *gorm.DB) error {
			// 再次检查锁
			if !lock1.CheckLockOwner(suite.Ctx) {
				return errors.New("锁已过期")
			}
			err := suite.productStockRepo.DeductStock(tx, productID, 1)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			suite.T().Logf("客户端1操作失败: %v", err)
		}
		// 尝试释放锁（即使已经过期）
		if err := lock1.Unlock(suite.Ctx); err != nil {
			suite.T().Logf("客户端1释放锁失败: %v", err)
		}
	}()

	// 模拟客户端2（保持不变）...

	// 模拟客户端3：等待锁过期后获取
	time.Sleep(5 * time.Second)

	go func() {
		defer wg.Done()
		time.Sleep(lockTimeout + time.Second)

		lock3 := NewRedisLock(suite.Rdb,
			fmt.Sprintf("product_lock_%d", productID),
			"client_3",
			lockTimeout,
		)

		ok, err := lock3.TryLock(suite.Ctx)
		suite.NoError(err)
		suite.True(ok, "锁已过期，客户端3应该能获取到锁")
		suite.T().Log("客户端3获取锁成功")

		defer lock3.Unlock(suite.Ctx)

		// 执行业务操作
		err = suite.Db.WithContext(suite.Ctx).Transaction(func(tx *gorm.DB) error {
			// 检查是否仍持有锁
			if !lock3.CheckLockOwner(suite.Ctx) {
				return errors.New("锁已失效")
			}

			err := suite.productStockRepo.DeductStock(tx, productID, 1)
			if err != nil {
				return err
			}
			suite.T().Log("客户端3执行操作成功")
			return nil
		})

		if err != nil {
			suite.T().Logf("客户端3操作失败: %v", err)
		}
	}()

	wg.Wait()
}

// 模拟三个用户同时下单，加入自旋等待获取锁机制
func (suite *RedisLockTestSuite) TestRedisLock_Scenario4() {
	var productID int64 = 123599
	lockTimeout := time.Second * 5

	// 自旋等待获取锁的最大次数
	maxRetries := 3
	// 自旋等待获取锁的间隔时间
	retryDelay := time.Second * 1

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func(userID int) {
			defer wg.Done()

			lock := NewRedisLock(suite.Rdb,
				fmt.Sprintf("product_lock_%d", productID),
				fmt.Sprintf("%d", userID),
				lockTimeout,
			)

			// 自旋等待锁
			var acquired bool
			var err error

			for retries := 0; retries < maxRetries; retries++ {
				acquired, err = lock.TryLock(suite.Ctx)
				// 获取锁失败，自旋继续获取
				if err != nil || !acquired {
					suite.T().Logf("用户%d第%d次获取锁失败，等待重试...", userID, retries+1)
					time.Sleep(retryDelay)
					continue
				}
				break
			}
			if !acquired {
				suite.T().Logf("用户%d在%d次尝试后仍未获取到锁，放弃下单", userID, maxRetries)
				return
			}
			suite.T().Logf("用户%d成功获取锁", userID)
			defer lock.Unlock(suite.Ctx)

			// 执行下单逻辑
			time.Sleep(1 * time.Second)
			// ...
		}(i + 1)
	}

	wg.Wait()
}

// 加入锁的续租
func (suite *RedisLockTestSuite) TestRedisLock_Scenario5() {
	var productID int64 = 123599
	lockTimeout := time.Second * 5

	var wg sync.WaitGroup
	wg.Add(2)

	// 用户1：获取锁并开启续租
	go func() {
		defer wg.Done()
		lock1 := NewRedisLock(suite.Rdb,
			fmt.Sprintf("product_lock_%d", productID),
			fmt.Sprintf("%d", 1),
			lockTimeout,
		)
		ok, err := lock1.TryLock(suite.Ctx)
		suite.NoError(err)
		suite.True(ok, "用户1成功获取到锁")
		// 开启锁自动续租
		lock1.EnableAutoRenew(suite.Ctx)
		defer func() {
			lock1.StopAutoRenew()
			lock1.Unlock(suite.Ctx)
		}()

		// 模拟长时间业务操作
		suite.T().Log("用户1开始执行长时间业务操作...")
		time.Sleep(lockTimeout * 2)

		err = suite.Db.WithContext(suite.Ctx).Transaction(func(tx *gorm.DB) error {
			if !lock1.CheckLockOwner(suite.Ctx) {
				return errors.New("锁已过期")
			}
			err := suite.productStockRepo.DeductStock(tx, productID, 1)
			if err != nil {
				return err
			}
			suite.T().Log("用户1扣减库存成功")
			return nil
		})
		suite.NoError(err, "由于自动续租，操作应该成功")
	}()

	// 用户2：在用户1执行期间尝试获取锁
	go func() {
		defer wg.Done()
		time.Sleep(lockTimeout + time.Second)
		lock2 := NewRedisLock(suite.Rdb,
			fmt.Sprintf("product_lock_%d", productID),
			fmt.Sprintf("%d", 2),
			lockTimeout,
		)

		ok, err := lock2.TryLock(suite.Ctx)
		suite.NoError(err)
		suite.False(ok, "由于用户1锁被续租，用户2应该无法获取到锁")
		suite.T().Log("用户2尝试获取锁失败(被用户1持有)")
	}()
	wg.Wait()
}
