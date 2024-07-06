package main

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/kafka_client"
	"gin-ddd-example/pkg/utils"

	"github.com/segmentio/kafka-go"
)

func init() {
	config.InitConfig()
	kafka_client.InitKafka(*config.Conf)
}

func main() {
	writterDemo2()
}

func simpleSendDemo() {
	topic := "test_simple_topic"
	partition := 0
	client := kafka_client.NewKafkaSimple(topic, partition)
	client.PublishSimple()
}

// 创建主题
func createTopicDemo() {
	topic := "user_registered"
	client := kafka_client.NewCreateTopic(topic, 2, 3)
	client.CreateTopic()
}

// 发送消息
func writterDemo() {
	client := kafka_client.NewKafkaWritter("user_registered")
	client.PublishMessage()
}

// 消息压缩示例
func writterDemo2() {
	client := kafka_client.NewKafkaWritter("test_simple_topic")
	client.PublishMessage2()
	fmt.Println("发送成功，继续后面的业务操作...")
}

// 获取 Topic 的分区信息
func getPartitionsDemo() {
	client := kafka_client.NewKafkaClient()
	conn, err := kafka.Dial("tcp", client.Dsn)
	client.FailOnErr(err, "Failed to dail broker")

	partitions, err := conn.ReadPartitions("user_registered")
	client.FailOnErr(err, "Failed to read partitions")

	utils.PrettyJson(partitions)
}
