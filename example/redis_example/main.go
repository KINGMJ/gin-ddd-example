package main

import (
	"context"
	"fmt"
	"gin-ddd-example/pkg/cache"
	"gin-ddd-example/pkg/config"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx context.Context

func main() {
	// 初始化操作
	config.InitConfig()
	cache.InitRedis(*config.Conf)
	rdb = cache.RedisClient
	ctx = context.Background()
	demo1()
}

func demo1() {
	val, err := rdb.Get(ctx, "vms_wx_access_token").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key:", val)
}
