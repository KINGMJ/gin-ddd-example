package cache_withlock_test

import (
	"fmt"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/service/cache_example/cache_withlock"
	"gin-ddd-example/pkg/test_suite"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type ProductServiceTestSuite struct {
	test_suite.TestSuite
	productRepo    repo.ProductRepo
	productService cache_withlock.ProductService
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.productRepo = repo.NewProductRepo()
	s.productService = cache_withlock.NewProductService(s.Ctx, s.Db, s.Rdb, s.Distlock, s.productRepo)
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
		product.Name = "冰淇淋"
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
