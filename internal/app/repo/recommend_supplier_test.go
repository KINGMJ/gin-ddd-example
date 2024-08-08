package repo_test

import (
	"encoding/json"
	"fmt"
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/utils"
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestRecommendFindById(t *testing.T) {
	// 必须定义为值类型，引用类型获取的为空
	var rs model.RecommendSupplierPo
	res := database.First(&rs, 1)

	fmt.Println("找到的记录数：", res.RowsAffected)
	fmt.Println("错误信息：", res.Error)
	utils.PrettyJson(rs)
}

func TestRecommendUpdate(t *testing.T) {
	// 使用 structs 更新
	logoData := make(map[string]any)
	logoData["url"] = "https://www.baiud.com"
	logoData["name"] = "哈哈11"
	// 将修改后的数据编码回 JSON
	updatedLogo, err := json.Marshal(logoData)
	if err != nil {
		log.Fatalf("Failed to marshal updated logo JSON: %v", err)
	}

	tx := database.Model(&model.RecommendSupplierPo{}).Where("id = ?", 1).
		Updates(model.RecommendSupplier{
			Logo: datatypes.JSON(updatedLogo),
		})
	fmt.Println(tx.Error)
}

func TestRecommendPreloadFindById(t *testing.T) {
	// 预加载关联关系
	var rs model.RecommendSupplierPo
	res := database.Preload("Type").Find(&rs, 1)

	fmt.Println("找到的记录数：", res.RowsAffected)
	fmt.Println("错误信息：", res.Error)
	utils.PrettyJson(rs)
}

func TestRecommendCreate(t *testing.T) {
	supplier := model.RecommendSupplier{
		Name:         gofakeit.Company(),
		Contact:      datatypes.JSON([]byte(`{"name":"张三","phone":"123456789"}`)),
		Logo:         datatypes.JSON([]byte(`{"url":"https://www.baiud.com","name":"哈哈"}`)),
		SpecialIntro: datatypes.JSON([]byte(`{"name":"哈哈","phone":"123456789"}`)),
		SpecialMark:  datatypes.JSON([]byte(`{"name":"哈哈","phone":"123456789"}`)),
	}

	supplierType := model.RecommendSupplierType{
		Name:   gofakeit.ProductCategory(),
		Status: 1,
	}

	rs := model.RecommendSupplierPo{
		RecommendSupplier: supplier,
		Type:              model.RecommendSupplierTypePo{RecommendSupplierType: supplierType},
	}
	tx := database.Create(&rs)
	fmt.Println(tx.Error)
	utils.PrettyJson(rs)
}

func TestRecommendSave(t *testing.T) {
	var rs model.RecommendSupplierPo
	database.First(&rs, 5)
	rs.Name = "供应商aaa"
	rs.TypeID = 1
	supplierType := model.RecommendSupplierType{
		Name:   "糖果",
		Status: 1,
	}
	rs.Type = model.RecommendSupplierTypePo{RecommendSupplierType: supplierType}
	tx := database.Save(&rs)
	fmt.Println(tx.Error)
	utils.PrettyJson(rs)
}

// 如果定义了关联关系，然后查询会怎么样
func TestRecommendSupplierBannerFindByID(t *testing.T) {
	// 必须定义为值类型，引用类型获取的为空
	var banner model.RecommendSupplierBannerPo
	res := database.Preload("RecommendSupplier").First(&banner, 1)

	fmt.Println("找到的记录数：", res.RowsAffected)
	fmt.Println("错误信息：", res.Error)
	utils.PrettyJson(banner)
}

// 多个 preload
func TestRecommendHasManyProducts(t *testing.T) {
	var rs model.RecommendSupplierPo
	res := database.Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Where(&model.RecommendSupplierProduct{TypeID: 1})
	}).Preload("Type").First(&rs, 1)
	fmt.Println("错误信息：", res.Error)
	utils.PrettyJson(rs)
}

// 列表里面preload是如何加载的
// Raw Sql:
// SELECT * FROM `recommend_supplier_product` WHERE `recommend_supplier_product`.`type_id` = 1 AND `recommend_supplier_product`.`supplier_id` IN (1,2,3)
func TestRecommendHasManyProducts2(t *testing.T) {
	var suppliers []*model.RecommendSupplierPo
	res := database.Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Where(&model.RecommendSupplierProduct{TypeID: 1})
	}).Limit(10).Find(&suppliers)
	fmt.Println("错误信息：", res.Error)
	utils.PrettyJson(suppliers)
}
