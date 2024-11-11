package repo

import (
	"errors"
	"fmt"
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/db"
	"gorm.io/gorm"
)

// 定义仓储接口
type SupplierRepo interface {
	FindById(id int) (*model.SupplierPo, error)
	Create(supplier *model.SupplierPo) (*model.SupplierPo, error)
}

type SupplierRepoImpl struct {
	*db.Database
}

func NewSupplierRepo(db *db.Database) *SupplierRepoImpl {
	return &SupplierRepoImpl{Database: db}
}

// ----------- (●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●) ------------

func (repo *SupplierRepoImpl) FindById(id int) (*model.SupplierPo, error) {
	// 必须定义为值类型，引用类型获取的为空
	var supplier model.SupplierPo
	// repo.DB.First(&supplier, id)
	res := repo.DB.Take(&supplier, id)

	// 根据指定字段查询，获取主键升序的第一条记录
	// res := repo.DB.First(&supplier, "merchant_id = ?", 15)

	fmt.Println("找到的记录数：", res.RowsAffected)
	fmt.Println("错误信息：", res.Error)

	// 判断查询的记录是否存在
	// if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }
	return &supplier, nil
}

func (repo *SupplierRepoImpl) FindByIds(ids []int) ([]*model.Supplier, error) {
	var suppliers []*model.Supplier
	// 查询所有的记录
	// res := repo.DB.Find(&suppliers)
	// 根据主键值检索
	res := repo.DB.Find(&suppliers, ids)
	// 判断查询的记录是否存在
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return suppliers, nil
}

func (repo *SupplierRepoImpl) FindByWhere() (model.Suppliers, error) {
	var suppliers model.Suppliers
	// 字符串条件，根据代码指定的顺序去查询
	// res := repo.DB.Where("name LIKE ?", "%有限公司%").
	// 	Where("merchant_id = ?", 15).
	// 	Find(&suppliers)

	// struct 条件，会按照索引查询
	// res := repo.DB.Where(&model.Supplier{
	// 	Scale:      1,
	// 	Mode:       1,
	// 	MerchantId: 15,
	// }).Find(&suppliers)

	// map 条件
	// res := repo.DB.Where(map[string]any{
	// 	"scale":       1,
	// 	"mode":        1,
	// 	"merchant_id": 15,
	// }).Find(&suppliers)

	// 内联条件
	// res := repo.DB.Find(&suppliers, &model.Supplier{
	// 	Scale:      1,
	// 	Mode:       1,
	// 	MerchantId: 15,
	// })

	// 多个条件
	// res := repo.DB.Where(&model.Supplier{MerchantId: 15}).
	// 	Select("name", "region").
	// 	Not(&model.Supplier{SType: 1}).
	// 	Not([]int64{1, 2, 3}).
	// 	Find(&suppliers)

	// 分页排序
	res := repo.DB.
		Order("created desc").
		Order("id asc").
		Limit(10).
		Offset(0).
		Find(&suppliers)

	// 判断查询的记录是否存在
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return suppliers, nil
}

func (repo *SupplierRepoImpl) Create(supplier *model.Supplier) (*model.Supplier, error) {
	res := repo.DB.Create(supplier)
	fmt.Println("影响的行数", res.RowsAffected)
	fmt.Println("错误信息：", res.Error)
	return supplier, res.Error
}

func (repo *SupplierRepoImpl) BatchCreate(suppliers model.Suppliers) (model.Suppliers, error) {
	// res := repo.DB.Create(suppliers)
	// 分批次执行
	res := repo.DB.CreateInBatches(suppliers, 3)
	fmt.Println("影响的行数", res.RowsAffected)
	fmt.Println("错误信息：", res.Error)
	return suppliers, res.Error
}

func (repo *SupplierRepoImpl) Save(supplier *model.Supplier) (*model.Supplier, error) {
	res := repo.DB.Save(supplier)
	fmt.Println("影响的行数", res.RowsAffected)
	fmt.Println("错误信息：", res.Error)
	return supplier, res.Error
}
