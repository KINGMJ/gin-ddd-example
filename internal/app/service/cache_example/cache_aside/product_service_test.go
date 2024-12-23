package cache_aside_test

import (
	"fmt"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/service/cache_example/cache_aside"
	"gin-ddd-example/pkg/test_suite"
	"gin-ddd-example/pkg/utils"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type ProductServiceTestSuite struct {
	test_suite.TestSuite
	productRepo    repo.ProductRepo
	productService cache_aside.ProductService
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.productRepo = repo.NewProductRepo()
	s.productService = cache_aside.NewProductService(s.Ctx, s.Db, s.Rdb, s.productRepo)
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

// 测试先删缓存，再更新db
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

// 测试先更新db，再删缓存
func (s *ProductServiceTestSuite) TestProductService_UpdateProductConcurrency2() {
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
		product.Name = "百世可乐"
		err := s.productService.UpdateProduct2(product)
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

// 测试延时双删
func (s *ProductServiceTestSuite) TestProductService_UpdateProductConcurrency3() {
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
		product.Name = "芒果"
		err := s.productService.UpdateProduct3(product)
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
