package db

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化db链接
func InitPostgresql() *Database {
	c := config.Conf.PostgresqlConf
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
		c.Host,
		c.User,
		c.Password,
		c.Dbname,
		c.Port,
		c.TimeZone,
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	Db = db
	return &Database{Db}
}
