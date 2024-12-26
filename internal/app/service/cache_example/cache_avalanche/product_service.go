package cache_avalanche

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gin-ddd-example/internal/app/model/product"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/container"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/sync/singleflight"
)

// 缓存雪崩的应对方案
type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*product.Product, error)
}

type productService struct {
	container.Container
	sf          *singleflight.Group
	productRepo repo.ProductRepo
}

func NewProductService(
	container container.Container,
	productRepo repo.ProductRepo,
) ProductService {
	return &productService{
		productRepo: productRepo,
		sf:          &singleflight.Group{},
		Container:   container,
	}
}

func (s *productService) GetProduct(ctx context.Context, id int64) (*product.Product, error) {
	start := time.Now()
	cacheKey := fmt.Sprintf("product:%d", id)
	// 记录请求信息
	var requestID string
	if id, ok := ctx.Value("request_id").(string); ok {
		requestID = id
	} else {
		requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
	}
	goroutineID := getGoroutineID() // 获取当前的 goroutine id

	s.Logs.Info(fmt.Sprintf("请求开始 - RequestID: %s, GoroutineID: %d, ProductID: %d", requestID, goroutineID, id))

	// 1. 先查询本地缓存
	productValue, err := s.getFromLocalCache(ctx, cacheKey)
	if err == nil {
		s.Logs.Info(fmt.Sprintf("本地缓存命中 - RequestID: %s, GoroutineID: %d", requestID, goroutineID))
		return productValue, nil
	}

	// 2. 缓存未命中，使用 single flight 合并请求结果
	value, err, shared := s.sf.Do(cacheKey, func() (interface{}, error) {
		// 记录实际执行查询的goroutine
		s.Logs.Info(fmt.Sprintf("执行数据库查询 - RequestID: %s, GoroutineID: %d", requestID, goroutineID))

		// 模拟数据库查询耗时
		time.Sleep(100 * time.Millisecond)

		product, err := s.productRepo.FindById(s.Db, id)
		if err != nil {
			return nil, err
		}

		// 异步更新缓存
		go func() {
			if err := s.updateCache(context.Background(), cacheKey, product); err != nil {
				s.Logs.Error(fmt.Sprintf("更新缓存失败: %v", err))
			}
		}()
		return product, nil
	})

	// 记录请求是否被合并
	if shared {
		s.Logs.Info(fmt.Sprintf("请求被合并 - RequestID: %s, GoroutineID: %d, 耗时: %v", requestID, goroutineID, time.Since(start)))
	} else {
		s.Logs.Info(fmt.Sprintf("独立请求 - RequestID: %s, GoroutineID: %d, 耗时: %v", requestID, goroutineID, time.Since(start)))
	}

	if err != nil {
		return nil, err
	}
	return value.(*product.Product), nil
}

func getGoroutineID() interface{} {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// getFromLocalCache 从本地缓存获取数据
func (s *productService) getFromLocalCache(ctx context.Context, key string) (*product.Product, error) {
	// 从 Redis 获取数据
	data, err := s.Rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var product product.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, err
	}
	return &product, nil
}

// updateCache 更新缓存
func (s *productService) updateCache(ctx context.Context, key string, product *product.Product) error {
	// 序列化数据
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	// 设置缓存，使用随机过期时间避免缓存雪崩
	expiration := getRandomExpiration()
	return s.Rdb.Set(ctx, key, data, expiration).Err()
}

// 过期时间随机化
func getRandomExpiration() time.Duration {
	baseExpiration := 30 * time.Minute
	// 随机添加0-5分钟
	random := time.Duration(rand.Int63n(300)) * time.Second
	return baseExpiration + random
}
