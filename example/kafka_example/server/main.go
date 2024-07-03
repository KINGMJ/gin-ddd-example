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
	writterDemo()
}

func simpleSendDemo() {
	topic := "test_simple_topic"
	partition := 0
	client := kafka_client.NewKafkaSimple(topic, partition)
	client.PublishSimple()
}

func createTopicDemo() {
	topic := "test_create_topic"
	client := kafka_client.NewCreateTopic(topic, 3, 3)
	client.CreateTopic()
}

func writterDemo() {
	client := kafka_client.NewKafkaWritter("test_create_topic")
	client.PublishMessage()
}
