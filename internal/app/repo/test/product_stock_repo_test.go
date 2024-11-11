package test

import (
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/utils"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProductStockRepoTestSuit struct {
	RepoTestSuite
	productStockRepo *repo.ProductStockRepoImpl
}

func (suite *ProductStockRepoTestSuit) SetupTest() {
	// 显示调用父类的 SetupTest，初始化 db
	suite.RepoTestSuite.SetupTest()
	suite.productStockRepo = repo.NewProductStockRepo(suite.db)
}

func TestProductStockRepoTestSuit(t *testing.T) {
	suite.Run(t, new(ProductStockRepoTestSuit))
}

func (suite *ProductStockRepoTestSuit) TestFindById() {
	res, err := suite.productStockRepo.FindById(1)
	if err != nil {
		suite.T().Error(err)
	}
	utils.PrettyJson(res)
}

func (suite *ProductStockRepoTestSuit) TestDeductStock() {
	err := suite.productStockRepo.DeductStock(suite.db.DB, 123599, 10)
	if err != nil {
		suite.T().Error(err)
	}
}
