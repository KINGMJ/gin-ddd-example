package redis_example

import (
	"context"
	"gin-ddd-example/pkg/cache"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	Rdb *redis.Client
	Db  *db.Database
	Ctx context.Context
}

func (suite *RedisTestSuite) SetupTest() {
	// 初始化操作
	config.InitConfig()
	// 日志初始化
	logs.InitLog(*config.Conf)
	cache.InitRedis(*config.Conf)
	suite.Ctx = context.Background()
	suite.Rdb = cache.RedisClient
	suite.Db = db.InitDb()
}
