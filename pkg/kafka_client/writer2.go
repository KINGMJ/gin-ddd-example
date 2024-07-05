package kafka_client

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func (client *KafkaClient) PublishMessage2() {
	// 连接 broker
	conn, err := kafka.Dial("tcp", client.Dsn)
	client.FailOnErr(err, "Failed to dail broker")
	defer conn.Close()

	// 读取 partitions
	partitions, err := conn.ReadPartitions(client.Topic)
	client.FailOnErr(err, "Failed to read partitions")

	// 获取 partitions 的 Leader 节点
	leadersMap := make(map[int]kafka.Broker)
	for _, item := range partitions {
		leadersMap[item.ID] = item.Leader
	}

	// 创建一个 writter 向 Topic 发送消息
	writter := kafka.Writer{
		Addr:        kafka.TCP(client.Dsn),
		Topic:       client.Topic,
		Balancer:    &kafka.RoundRobin{}, // 轮询策略
		Compression: kafka.Gzip,          // 使用 Gzip 压缩
	}

	var msgs []kafka.Message

	// 发送10条消息
	for i := 0; i < 10; i++ {
		msg := kafka.Message{
			Value: []byte(fmt.Sprintf("Hello Kafka %d", i)),
		}
		msgs = append(msgs, msg)
	}
	err = writter.WriteMessages(context.Background(), msgs...)
	client.FailOnErr(err, "Failed to write msg")
	writter.Close()
}
