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
	topicReceive()
}

func simpleReceive() {
	mq := rabbitmq.NewRabbitMQSimple("hello")
	mq.ReceiveSimple()
}

func workReceive() {
	mq := rabbitmq.NewRabbitMQSimple("work")
	mq.ReceiveWork()
}

func pubSubReceive() {
	mq := rabbitmq.NewRabbitMQPubSub("logs")
	mq.ReceiveSub()
}

func routingReceive() {
	mq := rabbitmq.NewRabbitMQRouting("logs_direct", "keyA")
	mq.ReceiveRouting()
}

func topicReceive() {
	mq := rabbitmq.NewRabbitMQTopic("logs_topic")
	mq.ReceiveTopic()
}
