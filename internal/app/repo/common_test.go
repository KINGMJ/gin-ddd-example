package repo_test

import (
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
)

var database *db.Database

func init() {
	config.InitConfig()
	// 日志初始化
	logs.InitLog(*config.Conf)
	logs.Log.Info("log init success!")
	database = db.InitDb()
}
