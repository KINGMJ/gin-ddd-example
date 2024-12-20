package service_test

import (
	"context"
	"fmt"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/service"
	"gin-ddd-example/pkg/cache"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"gin-ddd-example/pkg/rabbitmq"
	"gin-ddd-example/pkg/utils"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ProductServiceTestSuite struct {
	suite.Suite
	db                  *gorm.DB
	rdb                 *redis.Client
	ctx                 context.Context
	productRepo         repo.ProductRepo
	productService      service.ProductService
	productServiceCache service.ProductServiceCache
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) SetupSuite() {
	// 初始化操作
	config.InitConfig()
	cache.InitRedis(*config.Conf)
	rabbitmq.InitRabbitmq(*config.Conf)
	// 日志初始化
	logs.InitLog(*config.Conf)
	database := db.InitDb()
	s.ctx = context.Background()
	s.db = database.DB
	s.rdb = cache.RedisClient
	s.productRepo = repo.NewProductRepo()
	s.productService = service.NewProductService(s.ctx, s.db, s.rdb, s.productRepo)
	s.productServiceCache = service.NewProductServiceCache(s.ctx, s.db, s.rdb, s.productRepo)
}

func (s *ProductServiceTestSuite) TestProductService_GetProduct() {
	product, err := s.productService.GetProduct(123464)
	if err != nil {
		s.T().Error(err)
	}
	utils.PrettyJson(product)
}

func (s *ProductServiceTestSuite) TestProductService_UpdateProduct() {
	product, err := s.productService.GetProduct(123464)
	if err != nil {
		s.T().Error(err)
	}
	product.Name = "可口可乐"
	err = s.productService.UpdateProduct(product)
	if err != nil {
		s.T().Error(err)
	}
	utils.PrettyJson(product)
}

func (s *ProductServiceTestSuite) TestProductService_UpdateProductConcurrency() {
	// 预热缓存
	product, err := s.productService.GetProduct(123464)
	if err != nil {
		s.T().Error(err)
	}

	var wg sync.WaitGroup

	// 更新操作
	wg.Add(1)
	go func() {
		defer wg.Done()
		product.Name = "薯片"
		err := s.productService.UpdateProduct(product)
		s.NoError(err)
	}()

	// 并发读操作
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p, err := s.productService.GetProduct(123464)
			s.NoError(err)
			fmt.Println(p.Name)
		}()
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
}

func (s *ProductServiceTestSuite) TestProductService_UpdateProductConcurrency2() {
	// 预热缓存
	product, err := s.productServiceCache.Get(123464)
	if err != nil {
		s.T().Error(err)
	}

	var wg sync.WaitGroup

	// 更新操作
	wg.Add(1)
	go func() {
		defer wg.Done()
		product.Name = "薯片"
		err := s.productServiceCache.Set(product)
		s.NoError(err)
	}()

	// 并发读操作
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p, err := s.productServiceCache.Get(123464)
			s.NoError(err)
			fmt.Println(p.Name)
		}()
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
}
