package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var Conf = new(Config)

type Config struct {
	AppConf        AppConf      `mapstructure:"app"`        // app 配置
	PostgresqlConf DBConf       `mapstructure:"postgresql"` // postgresql数据库信息
	MysqlConf      DBConf       `mapstructure:"mysql"`      // mysql 数据库配置
	RedisConf      RedisConf    `mapstructure:"redis"`      // redis 配置
	RabbitmqConf   RabbitmqConf `mapstructure:"rabbitmq"`   // rabbitmq 配置
	LogsConf       LogsConf     `mapstructure:"logs"`       // 日志配置
	KafkaConf      KafkaConf    `mapstructure:"kafka"`      // kafka 配置
}

type AppConf struct {
	Env      string `mapstructure:"env"`
	Database string `mapstructure:"database"`
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

type KafkaConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type LogsConf struct {
	Level      string `mapstructure:"level"`       // 日志级别
	RootDir    string `mapstructure:"root_dir"`    // 日志文件存放位置
	Format     string `mapstructure:"format"`      // 格式：json 或者其他格式
	Filename   string `mapstructure:"filename"`    // 日志文件名
	DbFileName string `mapstructure:"db_filename"` // db 慢日志
	MaxSize    int    `mapstructure:"max_size"`    // 日志文件最大大小(M)
	MaxBackups int    `mapstructure:"max_backups"` // 旧文件的最大个数
	MaxAge     int    `mapstructure:"max_age"`     // 旧文件的最大保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
	ShowLine   bool   `mapstructure:"show_line"`   // 是否显示调用行
}

func InitConfig() {
	configPath := getConfigPath()
	// 指定配置文件路径
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// 读取配置信息
	err := viper.ReadInConfig()
	if err != nil { // 读取配置信息失败
		log.Fatal(fmt.Errorf("fatal error config file: %s", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatal(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}
	log.Println("The application configuration file is loaded successfully!")
}

func getConfigPath() string {
	workDir := os.Getenv("GIN_DDD_EXAMPLE")
	if workDir == "" {
		return "configs/dev.yaml"
	}
	return filepath.Join(workDir, "configs/dev.yaml")
}
