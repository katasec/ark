package messaging

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// RedisMessenger is an Azure Servicebus Implementation of the Messenger interface.
type RedisMessenger struct {
	ctx    context.Context
	client *redis.Client

	queueName string
}

func NewRedisMessenger(connectionString string, queueName string) *RedisMessenger {

	// Create context for Azure Service Bus
	ctx := context.Background()

	client := getRedisClient(connectionString)

	// Create ASB Messenger struct
	redisMessenger := &RedisMessenger{
		ctx:    ctx,
		client: client,

		queueName: queueName,
	}

	return redisMessenger
}

func (m *RedisMessenger) Send(subject string, message string) error {
	err := m.client.LPush(m.ctx, m.queueName, message).Err()
	return err
}

func (m *RedisMessenger) Receive() (string, string, error) {
	messageBody, err := m.client.LPop(m.ctx, m.queueName).Result()
	subject := m.queueName

	return messageBody, subject, err
}

func getRedisClient(connectionString string) *redis.Client {

	// Parse connection string
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Redis connection string was not in the correct format:" + connectionString)
	}

	// Create redis client
	client := redis.NewClient(opt)

	return client
}
