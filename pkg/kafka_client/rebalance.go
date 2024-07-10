package kafka_client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// rebalance 导致重复消费
func (client *KafkaClient) RebalanceReceive(id int, ctx context.Context, wg *sync.WaitGroup) {
	fmt.Printf("启动消费者实例：%d\n", id)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{client.Dsn},
		GroupID:        "group1",
		GroupTopics:    []string{"topic1"},
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 10 * time.Second,
	})
	defer reader.Close()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("退出接收...")
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				fmt.Printf("Failed to read message from topic %s: %v\n", msg.Topic, err)
				continue
			}
			fmt.Printf("Message from client %d %s at offset %d: %s = %s\n", id, msg.Topic, msg.Offset, msg.Key, msg.Value)
		}
	}
}
