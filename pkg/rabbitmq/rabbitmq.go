package rabbitmq

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"log"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
)

var MQURL string

func InitRabbitmq(config config.Config) {
	conf := config.RabbitmqConf
	// password 特殊字符转义
	var password = url.QueryEscape(conf.Password)
	// 连接信息
	MQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.User, password, conf.Host, conf.Port)
}

// rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// 断开channel 和 connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}
