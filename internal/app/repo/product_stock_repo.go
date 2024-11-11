package repo

import (
	"errors"
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/db"

	"gorm.io/gorm"
)

// ProductStockRepo 定义仓储接口
type ProductStockRepo interface {
	FindById(id int64) (*model.ProductStockPo, error)
	FindByProductId(productId int64) (*model.ProductStockPo, error)
	Create(stock *model.ProductStockPo) (*model.ProductStockPo, error)
	Update(stock *model.ProductStockPo) (*model.ProductStockPo, error)
	DeductStock(tx *gorm.DB, productID int64, count int64) error
}

type ProductStockRepoImpl struct {
	*db.Database
}

func NewProductStockRepo(db *db.Database) *ProductStockRepoImpl {
	return &ProductStockRepoImpl{Database: db}
}

func (repo *ProductStockRepoImpl) FindById(id int64) (*model.ProductStockPo, error) {
	var stock model.ProductStockPo
	res := repo.DB.Take(&stock, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &stock, res.Error
}

func (repo *ProductStockRepoImpl) FindByProductId(productId int64) (*model.ProductStockPo, error) {
	var stock model.ProductStockPo
	res := repo.DB.Where("product_id = ?", productId).Take(&stock)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &stock, res.Error
}

func (repo *ProductStockRepoImpl) Create(stock *model.ProductStockPo) (*model.ProductStockPo, error) {
	res := repo.DB.Create(stock)
	return stock, res.Error
}

func (repo *ProductStockRepoImpl) Update(stock *model.ProductStockPo) (*model.ProductStockPo, error) {
	res := repo.DB.Save(stock)
	return stock, res.Error
}

func (repo *ProductStockRepoImpl) DeductStock(tx *gorm.DB, productID int64, count int64) error {
	result := tx.Model(&model.ProductStockPo{}).
		Where("product_id = ?", productID).
		Update("count", gorm.Expr("count - ?", count))
	return result.Error
}
