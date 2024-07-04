package kafka_client

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/segmentio/kafka-go"
)

func NewKafkaWritter(topic string) *KafkaClient {
	client := NewKafkaClient()
	client.Topic = topic
	return client
}

func (client *KafkaClient) PublishMessage() {
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
		Addr:  kafka.TCP(client.Dsn),
		Topic: client.Topic,
		// Balancer: &kafka.RoundRobin{}, // 轮询策略
		// Balancer: &kafka.Hash{}, //根据 Key Hash 取模
		// Balancer: &kafka.LeastBytes{}, // 将消息发送到当前负载最小的分区
		// Balancer: &kafka.CRC32Balancer{}, // CRC32 校验
		// Balancer: &kafka.Murmur2Balancer{},
		// Balancer: &RandomBalancer{},
		// Balancer: &kafka.ReferenceHash{},
		Balancer: &GeoBalancer{Leaders: leadersMap}, // 自定义地理位置策略
	}

	var msgs []kafka.Message

	// 发送10条消息
	for i := 0; i < 10; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("South-%d", i)),
			Value: []byte(fmt.Sprintf("Hello Kafka %d", i)),
		}
		msgs = append(msgs, msg)
	}

	for i := 0; i < 10; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("North-%d", i)),
			Value: []byte(fmt.Sprintf("Hello Kafka2 %d", i)),
		}
		msgs = append(msgs, msg)
	}
	err = writter.WriteMessages(context.Background(), msgs...)
	client.FailOnErr(err, "Failed to write msg")
	writter.Close()
}

// 随机策略
type RandomBalancer struct{}

func (r *RandomBalancer) Balance(msg kafka.Message, partitions ...int) int {
	return rand.Intn(len(partitions))
}

// 地理位置策略
type GeoBalancer struct {
	Leaders map[int]kafka.Broker
}

func (r *GeoBalancer) Balance(msg kafka.Message, partitions ...int) int {
	// 这里判断是否是南方还是北方，基于Leader 的 IP 进行判断
	// 我们加上端口 9092是北方；9094是南方
	isNorth := strings.HasPrefix(string(msg.Key), "North")
	var northKey, southKey int
	for key, value := range r.Leaders {
		if value.Port == 9092 {
			northKey = key
		} else {
			southKey = key
		}
	}
	if isNorth {
		return northKey
	}
	return southKey
}
