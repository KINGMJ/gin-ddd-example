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
	simpleReceive()
}

func simpleReceive() {
	mq := rabbitmq.NewRabbitMQSimple("hello")
	mq.ReceiveSimple()
}
