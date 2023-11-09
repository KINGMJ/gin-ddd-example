package db

import (
	"fmt"
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 定义一个 database 结构体
type Database struct {
	DB *gorm.DB
}

// 初始化db链接
func InitDb() *Database {
	c := config.Conf.DBConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
		c.Host,
		c.User,
		c.Password,
		c.Dbname,
		c.Port,
		c.TimeZone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.User{}, &model.Ent{})
	return &Database{DB: db}
}
