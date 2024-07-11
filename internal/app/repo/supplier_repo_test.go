package repo_test

import (
	"fmt"
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/internal/app/model/ctype"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"gin-ddd-example/pkg/utils"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

var database *db.Database

func init() {
	config.InitConfig()
	// 日志初始化
	logs.InitLog(*config.Conf)
	logs.Log.Info("log init success!")
	database = db.InitDb()
}

func TestFindById(t *testing.T) {
	supplierRepo := &repo.SupplierRepoImpl{database}
	supplier, err := supplierRepo.FindById(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(supplier)
	utils.PrettyJson(supplier)
}

func TestFindByIds(t *testing.T) {
	supplierRepo := &repo.SupplierRepoImpl{database}
	suppliers, err := supplierRepo.FindByIds([]int{1, 2, 3})
	if err != nil {
		fmt.Println(err)
	}
	utils.PrettyJson(suppliers)
}

func TestFindByWhere(t *testing.T) {
	supplierRepo := &repo.SupplierRepoImpl{database}
	suppliers, err := supplierRepo.FindByWhere()
	if err != nil {
		fmt.Println(err)
	}
	utils.PrettyJson(suppliers)
}

func TestSave(t *testing.T) {
	supplierRepo := &repo.SupplierRepoImpl{database}
	res, err := supplierRepo.Create(&model.Supplier{
		Name:      gofakeit.Company(),
		SType:     1,
		Region:    gofakeit.City(),
		ComMobile: gofakeit.Phone(),
		Fax:       gofakeit.Phone(),
		BName:     gofakeit.Username(),
		BMobile:   gofakeit.Phone(),
		TaxesCard: gofakeit.CreditCard().Number,
		Created:   ctype.NewNullTime(time.Now()),
		Updated:   ctype.NewNullTime(time.Now()),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrettyJson(res)
}

func TestBatchSave(t *testing.T) {
	supplierRepo := &repo.SupplierRepoImpl{database}

	// 批量插入10条
	var suppliers model.Suppliers

	for i := 0; i < 10; i++ {
		suppliers = append(suppliers, &model.Supplier{
			Name:      gofakeit.Company(),
			SType:     1,
			Region:    gofakeit.City(),
			ComMobile: gofakeit.Phone(),
			Fax:       gofakeit.Phone(),
			BName:     gofakeit.Username(),
			BMobile:   gofakeit.Phone(),
			TaxesCard: gofakeit.CreditCard().Number,
			Created:   ctype.NewNullTime(time.Now()),
			Updated:   ctype.NewNullTime(time.Now()),
		})
	}

	res, err := supplierRepo.BatchCreate(suppliers)
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrettyJson(res)
}

// 更新，如果是部分更新，
// bank name 如何区分0值
// name= "bbb"  bName="" ,
// update  zhuan entity
