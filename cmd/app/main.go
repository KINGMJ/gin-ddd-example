package main

import (
	"gin-ddd-example/internal/app/model"
	"gin-ddd-example/pkg/db"
)

func main() {
	// // 初始化操作
	// config.InitConfig()
	// db.InitDb()
	// r := route.InitRouter()
	// // 运行服务
	// r.Run()

	db.Db.AutoMigrate(&model.Ent{})
}
