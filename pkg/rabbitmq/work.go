package rabbitmq

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 创建Work模式下 RabbitMQ 实例
func NewRabbitMQWork(queueName string) *RabbitMQ {
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
func (r *RabbitMQ) PublishWork() {
	// 1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		true,
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
	// 从命令中接收 body
	body := bodyFrom(os.Args)

	//调用 channel 发送消息到队列中
	err = r.channel.PublishWithContext(ctx,
		r.Exchange,  // 交换机
		r.QueueName, // routing key 名称，设置为队列名称
		false,       // mandatory， 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 表示消息需要持久化
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	r.failOnErr(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func (r *RabbitMQ) ReceiveWork() {
	// 1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	queue, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		true,
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

	// 设置消费者的服务质量（Quality of Service）
	// 可以帮助控制消息的分发，防止某个消费者被过多的消息压倒，而其他消费者却处于空闲状态
	err = r.channel.Qos(
		1,     // prefetch count 指定每个消费者预取的消息数量。例如，设置为 1 表示每个消费者一次从队列中获取一条消息。
		0,     // prefetch size 预取的消息大小，一般设置为 0 表示不限制消息大小。
		false, // global 如果为 true，表示设置的 prefetchCount 和 prefetchSize 应用于整个 channel，而不是每个消费者。
	)
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
			// 模拟耗时操作
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			// 表示拒绝消息，消息将被重新放回队列以便重新处理，需要将 auto-ack 设置为 false
			// 如果为 true，表示确认消息，即消费者成功处理了消息，消息可以从队列中移除。
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func bodyFrom(args []string) string {
	var s string
	// 如果没有传递命令行参数
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
