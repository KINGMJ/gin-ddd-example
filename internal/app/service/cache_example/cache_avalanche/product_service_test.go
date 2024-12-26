package cache_avalanche_test

import (
	"context"
	"fmt"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/service/cache_example/cache_avalanche"
	"gin-ddd-example/pkg/test_suite"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type ProductServiceTestSuite struct {
	test_suite.TestSuite
	productRepo    repo.ProductRepo
	productService cache_avalanche.ProductService
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.productRepo = repo.NewProductRepo()
	s.productService = cache_avalanche.NewProductService(s.Container, s.productRepo)
}

func (s *ProductServiceTestSuite) TestProductService_GetProductConcurrency() {
	var wg sync.WaitGroup
	// 并发读操作
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			// 为每个请求创建带有 request_id 的上下文
			ctx := context.WithValue(s.Ctx, "request_id", fmt.Sprintf("test-req-%d", index))
			p, err := s.productService.GetProduct(ctx, 123464)
			s.NoError(err)
			fmt.Println(p.Name)
		}(i)
	}
	wg.Wait()
}
