package kafka_client

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/segmentio/kafka-go"
)

// 模拟网络抖动发送消息是否会丢失

func (client *KafkaClient) PublishMessage2() {
	// 连接 broker
	conn, err := kafka.Dial("tcp", client.Dsn)
	client.FailOnErr(err, "Failed to dail broker")
	defer conn.Close()

	// 创建一个 writter 向 Topic 发送消息
	writter := kafka.Writer{
		Addr:        kafka.TCP(client.Dsn),
		Topic:       client.Topic,
		Balancer:    &kafka.RoundRobin{}, // 轮询策略
		Compression: kafka.Gzip,          // 使用 Gzip 压缩
	}

	var msgs []kafka.Message

	// 发送10条消息
	for i := 0; i < 10000; i++ {
		msg := kafka.Message{
			Value: []byte(fmt.Sprintf("Hello Kafka %d", i)),
		}
		msgs = append(msgs, msg)
	}

	// 模拟网络抖动，将网络设置为弱网模式
	resultChan := make(chan string)
	errorChan := make(chan error)
	go func() {
		cmd := exec.Command("sudo", "pfctl", "-e")
		// 获取命令的输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			errorChan <- err
		}
		resultChan <- string(output)
	}()

	err = writter.WriteMessages(context.Background(), msgs...)

	// 等待并接收协程的结果
	select {
	case output := <-resultChan:
		fmt.Printf("Command output: %s\n", output)
	case err := <-errorChan:
		fmt.Printf("Error: %s\n", err)
	}

	client.FailOnErr(err, "Failed to write msg")
	writter.Close()
}
