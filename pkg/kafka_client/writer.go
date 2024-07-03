package kafka_client

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/segmentio/kafka-go"
)

func NewKafkaWritter(topic string) *KafkaClient {
	client := NewKafkaClient()
	client.Topic = topic
	return client
}

func (client *KafkaClient) PublishMessage() {
	// 创建一个 writter 向 Topic 发送消息
	writter := kafka.Writer{
		Addr:  kafka.TCP(client.Dsn),
		Topic: client.Topic,
		// Balancer: &kafka.RoundRobin{}, // 轮询策略
		Balancer: &kafka.Hash{}, //根据 Key Hash 取模
		// Balancer: &kafka.LeastBytes{}, // 将消息发送到当前负载最小的分区
		// Balancer: &kafka.CRC32Balancer{}, // CRC32 校验
		// Balancer: &kafka.Murmur2Balancer{},
		// Balancer: &RandomBalancer{},
		// Balancer: &kafka.ReferenceHash{},
	}

	var msgs []kafka.Message
	// 发送10条消息
	for i := 0; i < 10; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", i)),
			Value: []byte(fmt.Sprintf("Hello Kafka %d", i)),
		}
		msgs = append(msgs, msg)
	}
	for i := 0; i < 10; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", i)),
			Value: []byte(fmt.Sprintf("Hello Kafka2 %d", i)),
		}
		msgs = append(msgs, msg)
	}

	err := writter.WriteMessages(context.Background(), msgs...)
	client.failOnErr(err, "Failed to write msg")

	writter.Close()
}

type RandomBalancer struct{}

func (r *RandomBalancer) Balance(msg kafka.Message, partitions ...int) int {
	return rand.Intn(len(partitions))
}
