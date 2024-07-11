package db

import (
	"fmt"
	"gin-ddd-example/internal/app/constants"
	"gin-ddd-example/pkg/config"
	"time"

	"github.com/natefinch/lumberjack"
	"gorm.io/gorm"
)

var Db *gorm.DB

// 定义一个 database 结构体
type Database struct {
	*gorm.DB
}

// 初始化db链接，根据配置文件选择连接的数据源
func InitDb() *Database {
	c := config.Conf
	if c.AppConf.Database == constants.DB_MYSQL {
		return InitMysql()
	}
	return InitPostgresql()
}

// 创建一个lumberjack.Logger对象，用于日志文件的切割和管理
func getLogWriter(config *config.Config) *lumberjack.Logger {
	// 按天进行日志分割
	today := time.Now().Format("2006-01-02")
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s-%s.log", config.LogsConf.RootDir, config.LogsConf.DbFileName, today),
		MaxSize:    config.LogsConf.MaxSize,
		MaxBackups: config.LogsConf.MaxBackups,
		MaxAge:     config.LogsConf.MaxAge,
		Compress:   config.LogsConf.Compress,
	}
}
