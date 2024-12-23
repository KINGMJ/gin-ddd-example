package write_through_example

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-ddd-example/internal/app/model/product"
	"gin-ddd-example/internal/app/repo"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Read/Write Through 模式
// 应用程序不直接与数据库交互，而是通过缓存层
type ProductServiceCache interface {
	Get(id int64) (*product.Product, error)
	Set(product *product.Product) error
}

type productServiceCache struct {
	ctx         context.Context
	db          *gorm.DB
	rdb         *redis.Client
	productRepo repo.ProductRepo
}

func NewProductServiceCache(
	ctx context.Context,
	db *gorm.DB,
	rdb *redis.Client,
	productRepo repo.ProductRepo,
) ProductServiceCache {
	return &productServiceCache{
		ctx:         ctx,
		db:          db,
		rdb:         rdb,
		productRepo: productRepo,
	}
}

// Read-Through 模式：缓存层负责从数据库加载数据
func (s *productServiceCache) Get(id int64) (*product.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)
	// 查询缓存
	data, err := s.rdb.Get(s.ctx, cacheKey).Result()
	if err == nil {
		var product product.Product
		err = json.Unmarshal([]byte(data), &product)
		return &product, err
	}
	// 缓存未命中，缓存层负责从数据库加载
	if !errors.Is(err, redis.Nil) {
		// 发生了其他错误
		return nil, err
	}
	product, err := s.productRepo.FindById(s.db, id)
	if err != nil {
		return nil, err
	}
	jsonData, _ := json.Marshal(product)
	err = s.rdb.Set(s.ctx, cacheKey, jsonData, 30*time.Minute).Err()
	return product, err
}

// Write-Through 模式：缓存层负责将数据写入数据库
func (s *productServiceCache) Set(product *product.Product) error {
	cacheKey := fmt.Sprintf("product:%d", product.ID)
	// 缓存层更新缓存
	jsonData, _ := json.Marshal(product)
	err := s.rdb.Set(s.ctx, cacheKey, jsonData, 30*time.Minute).Err()
	if err != nil {
		return err
	}
	// 模拟更新数据库延迟，线程B读取到旧值
	time.Sleep(1 * time.Second)
	return s.productRepo.Update(s.db, product)
}
