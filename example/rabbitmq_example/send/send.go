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
	simpleSend()
}

func simpleSend() {
	mq := rabbitmq.NewRabbitMQSimple("hello")
	mq.PublishSimple("使用简单模式发送的信息")
}
