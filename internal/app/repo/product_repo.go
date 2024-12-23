package repo

import (
	"gin-ddd-example/internal/app/model/product"

	"gorm.io/gorm"
)

type ProductRepo interface {
	FindById(tx *gorm.DB, id int64) (*product.Product, error)
	Update(tx *gorm.DB, product *product.Product) error
}

type productRepo struct{}

func NewProductRepo() ProductRepo {
	return &productRepo{}
}

func (repo *productRepo) FindById(tx *gorm.DB, id int64) (*product.Product, error) {
	var product product.Product
	res := tx.First(&product, id)
	return &product, res.Error
}

func (repo *productRepo) Update(tx *gorm.DB, product *product.Product) error {
	return tx.Save(product).Error
}
