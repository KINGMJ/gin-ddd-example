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
	simpleSend()
}

func simpleSend() {
	topic := "test_simple_topic"
	partition := 0
	client := kafka_client.NewKafkaSimple(topic, partition)
	client.PublishSimple()
}
