package repo_test

import (
	"fmt"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"gin-ddd-example/pkg/utils"
	"testing"
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
