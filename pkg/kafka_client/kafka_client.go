package kafka_client

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Conn      *kafka.Conn
	Topic     string
	Partition int
	// 连接信息
	Dsn string
}

var dsn string

func InitKafka(config config.Config) {
	dsn = fmt.Sprintf("%s:%s", config.KafkaConf.Host, config.KafkaConf.Port)
}

// 创建结构体实例
func NewKafkaClient(topic string, partition int) *KafkaClient {
	return &KafkaClient{Dsn: dsn, Topic: topic, Partition: partition}
}

// 断开连接
func (r *KafkaClient) Close() {
	err := r.Conn.Close()
	r.failOnErr(err, "Failed to close client")
}

// 错误处理函数
func (r *KafkaClient) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}
