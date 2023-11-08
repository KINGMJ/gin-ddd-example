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
