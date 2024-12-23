package cache_aside

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-ddd-example/internal/app/model/product"
	"gin-ddd-example/internal/app/repo"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 商品服务，用于测试 redis 的缓存

type ProductService interface {
	GetProduct(id int64) (*product.Product, error)
	UpdateProduct(product *product.Product) error
	UpdateProduct2(product *product.Product) error
	UpdateProduct3(product *product.Product) error
}

type productService struct {
	ctx         context.Context
	db          *gorm.DB
	rdb         *redis.Client
	productRepo repo.ProductRepo
}

func NewProductService(ctx context.Context, db *gorm.DB, rdb *redis.Client, productRepo repo.ProductRepo) ProductService {
	return &productService{
		ctx, db, rdb, productRepo,
	}
}

// -------------- Cache Aside 模式  --------------
func (s *productService) GetProduct(id int64) (*product.Product, error) {
	// 从 redis 中获取商品信息
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

	product, err := s.productRepo.FindById(s.db, id)
	if err != nil {
		return nil, err
	}

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
	// 模拟更新数据库延迟，线程B读取到旧值
	time.Sleep(1 * time.Second)
	// 更新数据库
	if err := s.productRepo.Update(s.db, product); err != nil {
		return err
	}
	return nil
}

// 先更新数据库，再删缓存
func (s *productService) UpdateProduct2(product *product.Product) error {
	// 更新数据库
	if err := s.productRepo.Update(s.db, product); err != nil {
		return err
	}
	// 模拟删除缓存延迟，线程B读取到旧值
	time.Sleep(100 * time.Millisecond)
	cacheKey := fmt.Sprintf("product:%d", product.ID)
	err := s.rdb.Del(s.ctx, cacheKey).Err()
	if err != nil {
		log.Printf("删除缓存失败: %v", err)
	}
	return nil
}
