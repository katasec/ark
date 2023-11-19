package messaging

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqMessenger struct {
	queue   amqp.Queue
	channel *amqp.Channel
	conn    *amqp.Connection
}

func NewRabbitMqMessenger(queueName string, connectionString string) *RabbitMqMessenger {

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	//defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return &RabbitMqMessenger{
		conn:    conn,
		queue:   q,
		channel: ch,
	}
}

func (msg *RabbitMqMessenger) Send(subject string, body string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := msg.channel.PublishWithContext(ctx,
		"",             // exchange
		msg.queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Type:        subject,
		})

	if err != nil {
		log.Printf("Error sending message: %s", err)
	}
	return err
}

func (msg *RabbitMqMessenger) Receive(subject string) <-chan amqp091.Delivery {

	msgs, err := msg.channel.Consume(
		msg.queue.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	return msgs
}

func (msg *RabbitMqMessenger) Close() {
	msg.channel.Close()
	msg.conn.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
