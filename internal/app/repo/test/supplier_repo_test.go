package test

import (
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/utils"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SupplierRepoTestSuit struct {
	RepoTestSuite
	supplierRepo *repo.SupplierRepoImpl
}

func (suite *SupplierRepoTestSuit) SetupTest() {
	// 显示调用父类的 SetupTest，初始化 db
	suite.RepoTestSuite.SetupTest()
	suite.supplierRepo = repo.NewSupplierRepo(suite.db)
}

func TestSupplierRepoTestSuit(t *testing.T) {
	suite.Run(t, new(SupplierRepoTestSuit))
}

func (suite *SupplierRepoTestSuit) TestFindByID() {
	supplier, err := suite.supplierRepo.FindById(6)
	if err != nil {
		suite.T().Error(err)
	}
	utils.PrettyJson(supplier)
}

// func TestFindByIds(t *testing.T) {
// 	supplierRepo := &repo.SupplierRepoImpl{database}
// 	suppliers, err := supplierRepo.FindByIds([]int{1, 2, 3})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	utils.PrettyJson(suppliers)
// }

// func TestFindByWhere(t *testing.T) {
// 	supplierRepo := &repo.SupplierRepoImpl{database}
// 	suppliers, err := supplierRepo.FindByWhere()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	utils.PrettyJson(suppliers)
// }

// func TestCreate(t *testing.T) {
// 	supplierRepo := &repo.SupplierRepoImpl{database}
// 	res, err := supplierRepo.Create(&model.Supplier{
// 		Name:      gofakeit.Company(),
// 		SType:     1,
// 		Region:    gofakeit.City(),
// 		ComMobile: gofakeit.Phone(),
// 		Fax:       gofakeit.Phone(),
// 		BName:     gofakeit.Username(),
// 		BMobile:   gofakeit.Phone(),
// 		TaxesCard: gofakeit.CreditCard().Number,
// 		// Created:   ctype.NewNullTime(time.Now()),
// 		// Updated:   ctype.NewNullTime(time.Now()),
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(res)
// 	utils.PrettyJson(res)
// }

// func TestBatchCreate(t *testing.T) {
// 	supplierRepo := &repo.SupplierRepoImpl{database}

// 	// 批量插入10条
// 	var suppliers model.Suppliers

// 	for i := 0; i < 10; i++ {
// 		suppliers = append(suppliers, &model.Supplier{
// 			Name:      gofakeit.Company(),
// 			SType:     1,
// 			Region:    gofakeit.City(),
// 			ComMobile: gofakeit.Phone(),
// 			Fax:       gofakeit.Phone(),
// 			BName:     gofakeit.Username(),
// 			BMobile:   gofakeit.Phone(),
// 			TaxesCard: gofakeit.CreditCard().Number,
// 			// Created:   ctype.NewNullTime(time.Now()),
// 			// Updated:   ctype.NewNullTime(time.Now()),
// 		})
// 	}

// 	res, err := supplierRepo.BatchCreate(suppliers)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	utils.PrettyJson(res)
// }

// // func TestSave(t *testing.T) {
// // 	supplierRepo := &repo.SupplierRepoImpl{database}
// // 	supplier, err := supplierRepo.FindById(1)
// // 	if err != nil {
// // 		fmt.Println(err)
// // 		return
// // 	}
// // 	supplier.Name = "bbb"
// // 	supplier.BName = ""
// // 	res, err := supplierRepo.Save(supplier)
// // 	if err != nil {
// // 		fmt.Println(err)
// // 		return
// // 	}
// // 	utils.PrettyJson(res)
// // }

// func TestUpdate(t *testing.T) {
// 	// 根据条件更新
// 	// tx := database.Model(&model.Supplier{}).Where("name = ?", "bbb").Update("name", "ccc")
// 	// 根据主键更新
// 	supplierRepo := &repo.SupplierRepoImpl{database}
// 	supplier, err := supplierRepo.FindById(2)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	tx := database.Model(supplier).Update("name", "hello")
// 	fmt.Println(tx.Error)
// }

// func TestUpdates(t *testing.T) {
// 	// 使用structs 更新
// 	// tx := database.Model(&model.Supplier{}).Where("s_type = ?", 1).Limit(1).
// 	// 	Updates(model.Supplier{Name: "jack", BName: "rose", Scale: 0})

// 	// 使用 map 更新
// 	tx := database.Model(&model.Supplier{}).Where("id = ?", 63).
// 		Updates(map[string]any{"name": "jack", "b_name": "rose", "scale": 0})
// 	fmt.Println(tx.Error)
// }

func (suite *SupplierRepoTestSuit) TestPreloadManyToMany() {
	// many to many 的 preload
	var supplier model.SupplierPo
	res := suite.db.Preload("Stores").First(&supplier, 2)
	if res.Error != nil {
		suite.T().Error(res.Error)
	}
	utils.PrettyJson(supplier)
}
