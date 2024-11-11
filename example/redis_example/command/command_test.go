package command

import (
	"fmt"
	"gin-ddd-example/example/redis_example"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	redis_example.RedisTestSuite
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) SetupTest() {
	suite.RedisTestSuite.SetupTest()
}

// --------------------------------------------------------------------
// 字符串操作测试
// --------------------------------------------------------------------
func (suite *CommandTestSuite) TestSet() {
	// set 如果 key 不存在，就设置；如果存在，就覆盖；可选过期时间
	// 使用场景：设置普通的缓存
	// err := suite.Rdb.Set(suite.Ctx, "name", "张三", 0).Err()
	// suite.NoError(err)

	// setnx 如果 key 不存在，就设置；如果存在，就返回 false
	// 使用场景：分布式锁、防止重复提交
	// ok, err := suite.Rdb.SetNX(suite.Ctx, "name", "李四", 0).Result()
	// suite.NoError(err)
	// suite.False(ok)

	// 专门用于设置会过期的key；必须设置过期时间，否则会报错
	// 使用场景：设置 token 、验证码等有过期时间的key
	err := suite.Rdb.SetEx(suite.Ctx, "name", "王五", time.Second*5).Err()
	suite.NoError(err)
}

// mset 批量设置，不支持过期时间
func (suite *CommandTestSuite) TestMSet() {
	// mset 同时设置多个 key-value
	// err := suite.Rdb.MSet(suite.Ctx, map[string]any{
	// 	"name": "赵六",
	// 	"age":  18,
	// }).Err()
	// suite.NoError(err)

	// msetnx 同时设置多个 key-value，如果有一个 key 存在，则全部设置失败
	ok, err := suite.Rdb.MSetNX(suite.Ctx, map[string]any{
		"name": "赵六",
		"age":  18,
	}).Result()
	suite.NoError(err)
	suite.False(ok)
}

func (suite *CommandTestSuite) TestGet() {
	// 1. 获取字符串，如果 key 不存在，则返回错误
	// val, err := suite.Rdb.Get(suite.Ctx, "name").Result()
	// suite.NoError(err)
	// suite.Equal("张三", val)

	// 2. GETSET - 设置新值并返回旧值，原子操作
	// oldVal, err := suite.Rdb.GetSet(suite.Ctx, "name", "李四").Result()
	// suite.NoError(err)
	// suite.Equal("张三", oldVal)

	// 3. MGET - 批量获取多个key的值
	// vals, err := suite.Rdb.MGet(suite.Ctx, "name", "age").Result()
	// suite.NoError(err)
	// suite.Equal([]any{"李四", "12"}, vals)

	// 4. GETRANGE - 获取字符串的子串。start 和 end 指的是字节位置（byte positon），而不是字符位置
	// start 和 end 都是包含在内的，可以使用负数表示从末尾开始计数
	// 对于中文，一个中文占 3 个字节，所以(0,2)才能表示一个中文字符
	val, err := suite.Rdb.GetRange(suite.Ctx, "name", -2, -1).Result()
	suite.NoError(err)
	fmt.Println(val)
}
