package container

import (
	"gin-ddd-example/pkg/cache"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/distlock"
	"gin-ddd-example/pkg/logs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	Db       *gorm.DB
	Rdb      *redis.Client
	Distlock distlock.Locker
	Logs     *zap.Logger
}

func NewContainer() Container {
	// 初始化操作
	config.InitConfig()
	cache.InitRedis(*config.Conf)
	// 日志初始化
	logs.InitLog(*config.Conf)
	database := db.InitDb()

	return Container{
		Db:       database.DB,
		Rdb:      cache.RedisClient,
		Distlock: distlock.NewDistLock(*config.Conf, cache.RedisClient, nil),
		Logs:     logs.Log,
	}
}
