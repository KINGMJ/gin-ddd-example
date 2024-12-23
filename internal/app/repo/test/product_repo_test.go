package test

import (
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/test_suite"
	"gin-ddd-example/pkg/utils"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProductRepoTestSuit struct {
	test_suite.TestSuite
	productRepo repo.ProductRepo
}

func (s *ProductRepoTestSuit) SetupTest() {
	s.TestSuite.SetupSuite()
	s.productRepo = repo.NewProductRepo()
}

func TestProductRepoTestSuit(t *testing.T) {
	suite.Run(t, new(ProductRepoTestSuit))
}

func (s *ProductRepoTestSuit) TestFindByID() {
	product, err := s.productRepo.FindById(s.Db, 123464)
	if err != nil {
		s.T().Error(err)
	}
	utils.PrettyJson(product)
}
