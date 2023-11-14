package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Conf = new(Config)

type Config struct {
	DBConf       DBConf       `mapstructure:"database"` // 数据库信息
	RedisConf    RedisConf    `mapstructure:"redis"`    // redis 配置
	RabbitmqConf RabbitmqConf `mapstructure:"rabbitmq"` // rabbitmq 配置
}

type DBConf struct {
	Dbname   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	TimeZone string `mapstructure:"timezone"`
}

type RedisConf struct {
	// redis服务器地址，ip:port格式，比如：192.168.1.100:6379
	// 默认为 :6379
	Addr     string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	// redis DB 数据库，默认为0
	Db int `mapstructure:"db"`
	// 连接池最大连接数量
	PoolSize int `mapstructure:"pool-size"`
	// 连接池保持的最小空闲连接数，它受到PoolSize的限制
	MinIdleConns int `mapstructure:"min-idle-conns"`
}

type RabbitmqConf struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func InitConfig() {
	var configPath string
	configEnv := os.Getenv("GO_ENV")
	switch configEnv {
	case "dev":
		configPath = "../../configs/dev.yml"
	case "test":
		configPath = "../../configs/test.yml"
	case "prod":
		configPath = "../../configs/prod.yml"
	default:
		// configPath = "../../configs/dev.yml"
		configPath = "configs/dev.yml"
	}
	// 指定配置文件路径
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// 读取配置信息
	err := viper.ReadInConfig()
	if err != nil { // 读取配置信息失败
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatal(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	log.Println("The application configuration file is loaded successfully!")
}
