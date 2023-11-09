package repo

import (
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/db"
)

// 定义 ent 仓储接口
type EntRepo interface {
	// Save(ent model.Ent) error
	// FindById(id int) (model.Ent, error)
	// Update(ent model.Ent) error
	List(page, pageSize int) ([]model.Ent, error)
}

type EntRepoImpl struct {
	db *db.Database
}

func NewEntRepo(db *db.Database) *EntRepoImpl {
	return &EntRepoImpl{db: db}
}

// ----------- (●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●) ------------

func (repo *EntRepoImpl) List(page, pageSize int) ([]model.Ent, error) {
	return []model.Ent{}, nil
}

// type entRepo struct {
// 	db *
// }

// 创建企业
// func Save(ent model.Ent) error {
// 	result := db.Db.Create(&ent)
// 	return result.Error
// }

// func FindById(id int) (model.Ent, error) {
// 	var ent model.Ent
// 	result := db.Db.First(&ent, id)
// 	return ent, result.Error
// }

// func Update(ent model.Ent) error {
// 	result := db.Db.Save(&ent)
// 	return result.Error
// }
