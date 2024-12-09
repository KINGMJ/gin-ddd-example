package repo

import (
	"gin-ddd-example/internal/app/model"
	"gorm.io/gorm"
)

type ProductRepo interface {
	FindById(tx *gorm.DB, id int) (*model.Product, error)
}

type productRepo struct{}

func NewProductRepo() ProductRepo {
	return &productRepo{}
}

func (repo *productRepo) FindById(tx *gorm.DB, id int) (*model.Product, error) {
	var product model.Product
	res := tx.First(&product, id)
	return &product, res.Error
}
