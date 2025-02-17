package cache_penetration_test

import (
	"context"
	"fmt"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/service/cache_example/cache_penetration"
	"gin-ddd-example/pkg/test_suite"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type ProductServiceTestSuite struct {
	test_suite.TestSuite
	productRepo    repo.ProductRepo
	productService cache_penetration.ProductService
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.productRepo = repo.NewProductRepo()
	s.productService = cache_penetration.NewProductService(s.Container, s.productRepo)
}

func (s *ProductServiceTestSuite) TestProductService_GetProductConcurrency() {
	var wg sync.WaitGroup
	results := make(chan error, 1000)

	// 并发读操作
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			// 为每个请求创建带有 request_id 的上下文
			ctx := context.WithValue(s.Ctx, "request_id", fmt.Sprintf("test-req-%d", index))
			p, err := s.productService.GetProduct(ctx, 123464)
			if err != nil {
				results <- err
				return // 有错误时直接返回，不继续执行
			}
			// 只有在没有错误时才打印
			fmt.Println(p.Name)
			results <- nil
		}(i)
	}

	wg.Wait()
	close(results)
	// 检查所有结果
	for err := range results {
		s.NoError(err)
	}
}
