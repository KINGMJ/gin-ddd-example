package rabbitmq

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	dlQueen      = "dead_letter_queue"
	dlExchange   = "dead_letter_exchange"
	dlRoutingKey = "dead_routing_key"
)

// 死信队列模式
func NewDead(exchangeName string, queenName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queenName, exchangeName, routingKey)
	var err error
	//获取 connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Failed to connect to Rabbitmq!")

	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Failed to open a channel")
	return rabbitmq
}

// 创建原始普通队列
func (r *RabbitMQ) CreateNormalQueue() {
	// 1. 创建原始普通队列，设置 ttl（模拟消息过期）、死信交换机、死信路由键
	_, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, amqp.Table{
		"x-message-ttl":             5000,         // 消息过期时间,毫秒
		"x-dead-letter-exchange":    dlExchange,   // 指定死信交换机
		"x-dead-letter-routing-key": dlRoutingKey, // 指定死信routing-key
	})

	if err != nil {
		r.failOnErr(err, "Failed to declare a queue")
	}
	// 2. 创建普通交换机
	err = r.channel.ExchangeDeclare(r.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		r.failOnErr(err, "Failed to declare a exchange")
	}
	// 3. 将队列、routing-key、交换机绑定
	err = r.channel.QueueBind(r.QueueName, r.Key, r.Exchange, false, nil)
	r.failOnErr(err, "Failed to bind a queue")
}

// 创建死信队列（队列、交换机、绑定）
func (r *RabbitMQ) CreateDeadQueue() {
	// 1. 创建死信队列
	_, err := r.channel.QueueDeclare(dlQueen, true, false, false, false, nil)
	if err != nil {
		r.failOnErr(err, "Failed to declare dead queue")
	}
	// 2. 创建死信交换机
	err = r.channel.ExchangeDeclare(dlExchange, "direct", true, false, false, false, nil)
	if err != nil {
		r.failOnErr(err, "Failed to declare dead exchange")
	}
	// 3. 将队列、routing-key、交换机绑定
	err = r.channel.QueueBind(dlQueen, dlRoutingKey, dlExchange, false, nil)
	r.failOnErr(err, "Failed to bind dead queue")
}

// 发布消息
func (r *RabbitMQ) PublishMessage() {
	r.CreateNormalQueue()
	r.CreateDeadQueue()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 从命令中接收 body
	body := bodyFrom(os.Args)

	// 2. 调用 channel 发送消息到队列中
	err := r.channel.PublishWithContext(ctx,
		r.Exchange, // 交换机
		r.Key,      // routing key 名称
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

// 接收消息，模拟丢弃消息
func (r *RabbitMQ) ReceiveMessage() {
	r.CreateNormalQueue()
	// 注册一个客户端
	msgs, err := r.channel.Consume(r.QueueName, "", false, false, false, false, nil)
	r.failOnErr(err, "Failed to register a consumer")
	var forever chan struct{}
	//启用协程处理消息
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			// 丢弃信息
			d.Nack(true, true)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
