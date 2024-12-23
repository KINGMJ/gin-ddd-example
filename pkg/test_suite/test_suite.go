package test_suite

import (
	"context"
	"gin-ddd-example/pkg/cache"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	Db  *gorm.DB
	Rdb *redis.Client
	Ctx context.Context
}

func (s *TestSuite) SetupSuite() {
	// 初始化操作
	config.InitConfig()
	cache.InitRedis(*config.Conf)
	// 日志初始化
	logs.InitLog(*config.Conf)
	database := db.InitDb()
	s.Ctx = context.Background()
	s.Db = database.DB
	s.Rdb = cache.RedisClient
}
