package messaging

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqMessenger struct {
	queue   amqp.Queue
	channel *amqp.Channel
	conn    *amqp.Connection
}

func NewRabbitMqMessenger(connectionString string, queueName string) *RabbitMqMessenger {

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// Configure QoS to limit unacknowledged messages to one
	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	// Declare/Create a queue
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

	ctx := context.Background()

	err := msg.channel.PublishWithContext(
		ctx,
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

func (msg *RabbitMqMessenger) ReceiveChannel(subject string) <-chan amqp091.Delivery {

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

func (msg *RabbitMqMessenger) Receive() (message string, subject string, err error) {

	msgs, err := msg.channel.Consume(msg.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	data := <-msgs

	return string(data.Body), string(data.Type), err
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
