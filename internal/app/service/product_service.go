package service

import (
	"context"
	"gin-ddd-example/internal/app/repo"
	"gorm.io/gorm"
)

// 商品服务，用于测试 redis 的缓存

type ProductService interface {
}

type productService struct {
	ctx         context.Context
	db          *gorm.DB
	productRepo repo.ProductRepo
}

func NewProductService(ctx context.Context, db *gorm.DB, productRepo repo.ProductRepo) ProductService {
	return &productService{
		ctx, db, productRepo,
	}
}
