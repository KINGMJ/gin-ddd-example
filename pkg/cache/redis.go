package cache

import (
	"context"
	"gin-ddd-example/pkg/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化redisClient
func InitRedis(config config.Config) {
	conf := config.RedisConf
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		DB:           conf.Db,
		Password:     conf.Password,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
	})
	// 心跳检测：使用Ping方法检查连接是否正常
	_, err := RedisClient.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
}
