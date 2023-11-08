package ent_dao

import (
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/db"
)

// 创建企业
func Save(ent model.Ent) error {
	result := db.Db.Create(&ent)
	return result.Error
}

func FindById(id int) (model.Ent, error) {
	var ent model.Ent
	result := db.Db.First(&ent, id)
	return ent, result.Error
}

func Update(ent model.Ent) error {
	result := db.Db.Save(&ent)
	return result.Error
}
