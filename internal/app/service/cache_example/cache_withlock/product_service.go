package cache_withlock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-ddd-example/internal/app/model/product"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/distlock"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 商品服务，用于测试 redis 的缓存

type ProductService interface {
	GetProduct(id int64) (*product.Product, error)
	UpdateProduct(product *product.Product) error
}

type productService struct {
	ctx         context.Context
	db          *gorm.DB
	rdb         *redis.Client
	distlock    distlock.Locker
	productRepo repo.ProductRepo
}

func NewProductService(
	ctx context.Context,
	db *gorm.DB,
	rdb *redis.Client,
	distlock distlock.Locker,
	productRepo repo.ProductRepo,
) ProductService {
	return &productService{
		ctx, db, rdb, distlock, productRepo,
	}
}

func (s *productService) GetProduct(id int64) (*product.Product, error) {
	// 1. 先读缓存
	cacheKey := fmt.Sprintf("product:%d", id)
	data, err := s.rdb.Get(s.ctx, cacheKey).Result()
	// 缓存命中
	if err == nil {
		var product product.Product
		err = json.Unmarshal([]byte(data), &product)
		return &product, err
	}
	if !errors.Is(err, redis.Nil) {
		// 发生了其他错误
		return nil, err
	}
	// 2. 缓存未命中，获取锁
	lock, err := s.distlock.Lock(s.ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	defer lock.Unlock(s.ctx)

	// 3. 双重检查（Double Check）
	data, err = s.rdb.Get(s.ctx, cacheKey).Result()
	// 缓存命中
	if err == nil {
		var product product.Product
		err = json.Unmarshal([]byte(data), &product)
		return &product, err
	}
	if !errors.Is(err, redis.Nil) {
		// 发生了其他错误
		return nil, err
	}

	// 4. 从数据库读取
	product, err := s.productRepo.FindById(s.db, id)
	if err != nil {
		return nil, err
	}
	// 5. 写入缓存
	jsonData, _ := json.Marshal(product)
	err = s.rdb.Set(s.ctx, cacheKey, jsonData, 30*time.Minute).Err()
	return product, err
}

// 先删缓存，再更新数据库
func (s *productService) UpdateProduct(product *product.Product) error {
	cacheKey := fmt.Sprintf("product:%d", product.ID)
	err := s.rdb.Del(s.ctx, cacheKey).Err()
	if err != nil {
		log.Printf("删除缓存失败: %v", err)
	}
	// 2. 获取锁
	lock, err := s.distlock.Lock(s.ctx, cacheKey)
	if err != nil {
		return err
	}
	defer lock.Unlock(s.ctx)

	// 模拟更新数据库延迟，线程B读取到旧值
	time.Sleep(1 * time.Second)
	// 更新数据库
	if err := s.productRepo.Update(s.db, product); err != nil {
		return err
	}
	return nil
}
