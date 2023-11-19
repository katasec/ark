package messaging

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/katasec/ark/config"
	"github.com/stretchr/testify/assert"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestMqConfig(t *testing.T) {
	config := config.ReadConfig()

	assert.NotEmpty(t, config.MqConnectionString, "The MqConnectionString was nil")
}

func TestSendMessage(t *testing.T) {
	config := config.ReadConfig()
	assert.NotEmpty(t, config.MqConnectionString, "The MqConnectionString was nil")

	conn, err := amqp.Dial(config.MqConnectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

// func TestSendMessage2(t *testing.T) {
// 	config := config.ReadConfig()
// 	assert.NotEmpty(t, config.MqConnectionString, "The MqConnectionString was nil")

// 	var mq = NewRabbitMqMessenger("commandqueue", config.MqConnectionString)
// 	mq.Send("test", "Hello World111!")

// 	message := mq.Receive("test")
// 	log.Printf("The received message was:" + message)
// }
