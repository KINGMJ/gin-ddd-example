package main

import (
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/kafka_client"
)

func init() {
	config.InitConfig()
	kafka_client.InitKafka(*config.Conf)
}

func main() {
	readerGroupDemo()
}

func simpleReceive() {
	topic := "test_simple_topic"
	partition := 0
	client := kafka_client.NewKafkaSimple(topic, partition)
	client.ReceiveSimple2()
}

// 模拟读取消息
func readerDemo() {
	client := kafka_client.NewKafkaReader("test_simple_topic")
	go client.ReceiveMessage3(0)
	go client.ReceiveMessage3(1)
	select {}
}

// 消费者组示例
func readerGroupDemo() {
	client := kafka_client.NewKafkaClient()
	client.ReceiveMessage4()
}
