package service_test

import (
	"context"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/service"
	"gin-ddd-example/pkg/cache"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"gin-ddd-example/pkg/rabbitmq"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type ProductServiceTestSuite struct {
	suite.Suite
	db             *gorm.DB
	ctx            context.Context
	productRepo    repo.ProductRepo
	productService service.ProductService
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
	s.productRepo = repo.NewProductRepo()
	s.productService = service.NewProductService(s.ctx, s.db, s.productRepo)
}

func (s *ProductServiceTestSuite) TestProductService_CreateProduct() {

}
