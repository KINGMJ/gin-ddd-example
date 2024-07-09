package kafka_client

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func NewKafkaReader(topic string) *KafkaClient {
	client := NewKafkaClient()
	client.Topic = topic
	return client
}

// 正常消费消息，每次重启都会重新消费分区内所有的消息
func (client *KafkaClient) ReceiveMessage3(partition int) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{client.Dsn},
		Topic:     client.Topic,
		MinBytes:  10e3,      // 10kb
		MaxBytes:  10e6,      // 10mb
		Partition: partition, // 指定消费哪个分区
	})
	defer reader.Close()

	// 消费消息
	for {
		msg, err := reader.ReadMessage(context.Background())
		client.FailOnErr(err, "Failed to read message")
		fmt.Printf("Message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
	}
}

// 手动管理偏移量
func (client *KafkaClient) ReceiveMessage3_1() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{client.Dsn},
		Topic:    client.Topic,
		MinBytes: 10e3, // 10kb
		MaxBytes: 10e6, // 10mb
	})
	defer reader.Close()

	// 从外部存储中获取上次处理的偏移量
	offset, err := getStoredOffset(client.Topic)
	client.FailOnErr(err, "Failed to get offset from store")

	// 设置偏移量
	reader.SetOffset(offset)

	// 消费消息
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Failed to read message: %v\n", err)
			// 在这里可以增加重试逻辑
			// 例如，如果是临时性错误，可以增加一个重试次数限制
			// 也可以记录到日志或者采取其他处理方式
			continue
		}
		// 其他业务操作...
		fmt.Printf("Message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
		// 处理完后记录当前偏移量到外部存储
		err = storeOffsetToStorage(client.Topic, msg.Offset+1)
		if err != nil {
			fmt.Println("Error storing offset to external storage:", err)
		}
	}
}

// 从外部存储获取偏移量
func getStoredOffset(topic string) (int64, error) {
	// 实现根据 topic 获取存储的偏移量逻辑，比如从数据库中查询
	// 这里简化为直接返回指定的值
	return 30005, nil
}

// 将偏移量存储到外部存储
func storeOffsetToStorage(topic string, offset int64) error {
	// 实现将偏移量存储到外部存储的逻辑，比如存储到数据库
	// 这里简化为打印偏移量
	fmt.Printf("Storing offset %d to external storage\n", offset)
	return nil
}
