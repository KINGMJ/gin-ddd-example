package main

import (
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/rabbitmq"
)

func init() {
	config.InitConfig()
	rabbitmq.InitRabbitmq(*config.Conf)
}

func main() {
	topicSend()
}

func simpleSend() {
	mq := rabbitmq.NewRabbitMQSimple("hello")
	mq.PublishSimple("使用简单模式发送的信息")
}

func workSend() {
	mq := rabbitmq.NewRabbitMQSimple("work")
	mq.PublishWork()
}

func pubSubSend() {
	mq := rabbitmq.NewRabbitMQPubSub("logs")
	mq.PublishPub()
}

func routingSend() {
	mq := rabbitmq.NewRabbitMQRouting("logs_direct", "keyA")
	mq.PublishRouting()
}

func topicSend() {
	mq := rabbitmq.NewRabbitMQTopic("logs_topic")
	mq.PublishTopic()
}
