package kafka_client

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

func (client *KafkaClient) ReceiveMessage4() {
	groupId := "example-group"
	// 订阅三个主题
	topics := []string{"board-events", "file-events", "task-events"}

	ctx := context.Background()

	// 消费者组的数量，最好设置为订阅的主题总分区数量
	for _, topic := range topics {
		// 假设每个Topic都有2个分区
		for i := 0; i < 2; i++ {
			go consumeTopic(ctx, client.Dsn, groupId, topic)
		}
	}
	// 阻塞主线程
	select {}
}

func consumeTopic(ctx context.Context, broker string, groupId, topic string) {
	// 创建 kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		// CommitInterval: time.Second, // 1秒间隔提交偏移量
	})
	defer reader.Close()
	// 消费消息
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Failed to read message from topic %s: %v\n", topic, err)
			// 消息失败处理
			// 如果是网络抖动导致读取不到，可以采用重试机制
			// 可以在这里记录日志，或者将消息发送到死信队列
			continue
		}
		fmt.Printf("Message from topic %s at offset %d: %s = %s\n", topic, msg.Offset, msg.Key, msg.Value)
	}
}

// 使用context来管理多个消费者进行消费
func (client *KafkaClient) ReceiveMessage5(ctx context.Context, wg *sync.WaitGroup) {
	// 假设Topic的分区数都是1，我们设置1个客户端
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{client.Dsn},
		GroupID:     "group1",
		GroupTopics: []string{"topic1", "topic2"},
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})
	defer reader.Close()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("退出接收")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				fmt.Printf("Failed to read message from topic %s: %v\n", msg.Topic, err)
				continue
			}
			fmt.Printf("Message from topic %s at offset %d: %s = %s\n", msg.Topic, msg.Offset, msg.Key, msg.Value)
		}
	}
}
