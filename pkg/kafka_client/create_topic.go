package kafka_client

import (
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// 创建一个topic

func NewCreateTopic(topic string, numPartitions, replicationFactor int) *KafkaClient {
	client := NewKafkaClient()
	// 连接至kafka任意节点
	conn, err := kafka.Dial("tcp", client.Dsn)
	client.FailOnErr(err, "Failed to dial broker")

	client.Conn = conn
	client.Topic = topic
	client.NumPartitions = numPartitions
	client.ReplicationFactor = replicationFactor

	return client
}

func (client *KafkaClient) CreateTopic() {
	// 获取当前控制节点信息
	// 控制器节点是Kafka集群中负责管理分区的领导选举和其他重要元数据操作的节点
	controller, err := client.Conn.Controller()
	client.FailOnErr(err, "Failed to get controller")

	// 连接至Leader节点
	conn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	client.FailOnErr(err, "Failed to dial leader")
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             client.Topic,
			NumPartitions:     client.NumPartitions,     // 分区数量
			ReplicationFactor: client.ReplicationFactor, // 每个分区的副本数，不能超过broker的数量，否则创建失败
		},
	}
	// 创建topic
	err = conn.CreateTopics(topicConfigs...)
	client.FailOnErr(err, "Failed to create topic")

	client.Close()
}
