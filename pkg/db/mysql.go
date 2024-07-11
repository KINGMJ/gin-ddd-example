package db

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"log"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化db链接
func InitMysql() *Database {
	c := config.Conf.MysqlConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Dbname,
		url.QueryEscape(c.TimeZone), // Asia/Shanghai，斜杠为特殊符号，需要转为%2F，否则报错：invalid DSN: did you forget to escape a param value?
	)

	newLogger := logger.New(
		log.New(getLogWriter(config.Conf), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * 100, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,                  // Don't include params in the SQL log
			Colorful:                  false,                  // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	Db = db
	return &Database{Db}
}
