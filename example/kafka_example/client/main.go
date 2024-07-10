package main

import (
	"context"
	"fmt"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/kafka_client"
	"os"
	"os/signal"
	"sync"
)

func init() {
	config.InitConfig()
	kafka_client.InitKafka(*config.Conf)
}

func main() {
	rebalanceDemo1()
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

func readerGroupDemo2() {
	client := kafka_client.NewKafkaClient()
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	// go client.ReceiveMessage4_1(ctx, &wg)
	go client.ReceiveMessage4_2(ctx, &wg)

	// 捕获中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Shutting down consumers...")
	cancel()
	wg.Wait()
}

// 重平衡示例
func rebalanceDemo1() {
	client := kafka_client.NewKafkaClient()
	var wg sync.WaitGroup
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())

	go client.ReceiveMessage4_1(ctx, &wg)
	go client.ReceiveMessage4_1(ctx, &wg)

	// 捕获中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Shutting down consumers...")
	cancel()
	wg.Wait()
}
