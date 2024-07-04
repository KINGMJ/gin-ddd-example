package kafka_client

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Conn              *kafka.Conn
	Dsn               string // 连接信息
	Topic             string // 主题
	NumPartitions     int    // 分区数量
	ReplicationFactor int    // 副本数量
}

var dsn string

func InitKafka(config config.Config) {
	dsn = fmt.Sprintf("%s:%s", config.KafkaConf.Host, config.KafkaConf.Port)
}

// 创建结构体实例
func NewKafkaClient() *KafkaClient {
	return &KafkaClient{Dsn: dsn}
}

// 断开连接
func (r *KafkaClient) Close() {
	err := r.Conn.Close()
	r.FailOnErr(err, "Failed to close client")
}

// 错误处理函数
func (r *KafkaClient) FailOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}
