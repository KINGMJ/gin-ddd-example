package rabbitmq

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 创建发布/订阅模式下 RabbitMQ 实例
// 使用交换机而不是队列名来传递消息
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	//获取 connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Failed to connect to Rabbitmq!")

	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Failed to open a channel")

	return rabbitmq
}

// 发布订阅模式下发送消息
func (r *RabbitMQ) PublishPub() {
	// 1. 创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, // 交换机名称
		"fanout",   // 交换机类型：扇出类型
		true,       // 持久
		false,      // 自动删除
		false,      // 内部，true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,      // 不等待
		nil,        // 额外参数
	)
	r.failOnErr(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 从命令中接收 body
	body := bodyFrom(os.Args)

	// 2. 调用 channel 发送消息到队列中
	err = r.channel.PublishWithContext(ctx,
		r.Exchange, // 交换机
		"",         // 队列名称
		false,      // mandatory， 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 表示消息需要持久化
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	r.failOnErr(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

// 接收消息
func (r *RabbitMQ) ReceiveSub() {
	// 1. 定义交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2.申请队列，如果队列不存在会自动创建，存在则跳过创建
	queue, err := r.channel.QueueDeclare(
		"",    // 创建一个随机队列
		true,  //是否持久化
		false, //是否自动删除
		false, //是否具有排他性
		false, //是否阻塞处理
		nil,   //额外的属性
	)
	r.failOnErr(err, "Failed to declare a queue")

	// 3. 将队列绑定到交换机上
	err = r.channel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		r.Exchange, // exchange
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind a queue")

	// 4. 注册一个客户端
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
			log.Printf(" [x] %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
