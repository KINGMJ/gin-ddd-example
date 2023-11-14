package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 创建简单模式下 RabbitMQ 实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	//获取 connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Failed to connect to Rabbitmq!")

	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Failed to open a channel")

	return rabbitmq
}

// 简单模式下发送消息
func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//调用 channel 发送消息到队列中
	err = r.channel.PublishWithContext(ctx,
		r.Exchange,  // 交换机
		r.QueueName, // 队列名称
		false,       // mandatory， 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", message)
}

func (r *RabbitMQ) ReceiveSimple() {
	// 1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	queue, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	// 注册一个客户端
	msgs, err := r.channel.Consume(
		queue.Name, // queue
		"",         // consumer, 用来区分多个消费者
		true,       // auto-ack, 是否自动应答
		false,      // exclusive,是否独有
		false,      // no-local, 设置为true，表示 不能将同一个Connection中生产者发送的消息传递给这个Connection中 的消费者
		false,      // no-wait, 队列是否阻塞
		nil,        // args, 额外的属性
	)
	r.failOnErr(err, "Failed to register a consumer")

	var forever chan struct{}
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
