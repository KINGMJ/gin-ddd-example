package kafka_client

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// 一个最简单的生产者和消费者示例
// 生产者和消费者都使用DialLeader连接到Leader Broker的指定分区
func NewKafkaSimple(topic string, partition int) *KafkaClient {
	client := NewKafkaClient()
	// 连接至 Kafka集群的Leader节点
	conn, err := kafka.DialLeader(context.Background(), "tcp", client.Dsn, topic, partition)
	client.failOnErr(err, "Failed to dial leader")
	client.Conn = conn
	return client
}

func (client *KafkaClient) PublishSimple() {
	// 设置发送消息的超时时间
	client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	// 发送消息
	_, err := client.Conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("tow!")},
		kafka.Message{Value: []byte("three!")},
	)
	client.failOnErr(err, "Failed to send message")
	// 关闭连接
	client.Close()
}

func (client *KafkaClient) ReceiveSimple() {
	// 设置读取超时时间
	client.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	// 读取一批消息，得到的batch是一系列消息的迭代器
	// 【注意：】这里最小值和最大值设置会影响 batch 的关闭。我们最小设置为 1e1，也就是10byte，
	batch := client.Conn.ReadBatch(1e1, 10e3) // e3 = 10^3(kb), e6= 10^6(mb)
	// 遍历读取消息
	b := make([]byte, 10e3) // 1kb max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Printf("[*] Receive %s\n", string(b[:n]))
	}
	// 关闭batch
	err := batch.Close()
	client.failOnErr(err, "Failed to close batch")
	// 关闭连接
	defer client.Close()
}

// 使用 batch.ReadMessage 读取消息
func (client *KafkaClient) ReceiveSimple2() {
	client.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	batch := client.Conn.ReadBatch(1e1, 10e3) // e3 = 10^3(kb), e6= 10^6(mb)
	for {
		msg, err := batch.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(string(msg.Value))
	}
	// 关闭batch
	err := batch.Close()
	client.failOnErr(err, "Failed to close batch")
	// 关闭连接
	defer client.Close()
}
