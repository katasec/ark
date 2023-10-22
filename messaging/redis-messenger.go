package messaging

import (
	"context"
	"log"
	"net/url"
	"strings"

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

func (m *RedisMessenger) Receive(subject string, message string) (string, string, error) {
	messageBody, err := m.client.LPop(m.ctx, m.queueName).Result()
	subject = m.queueName

	return messageBody, subject, err
}

func getRedisClient(connectionString string) *redis.Client {

	// validate connection string
	//connectionString := "redis://localhost:6379"
	url, err := url.Parse(connectionString)
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Redis connection string was not in the correct format:" + connectionString)
	}

	// get redis config from connection string
	config := getRedisConfig(url)

	// Connect to Redis.
	client := redis.NewClient(&redis.Options{
		Addr: config["server"],
	})

	//return client
	return client
}

func getRedisConfig(u *url.URL) map[string]string {
	config := map[string]string{
		"server":   u.Host,
		"database": strings.Trim(u.Path, "/"),
		"pool":     "30",
		"process":  "1",
	}

	for k, v := range u.Query() {
		config[k] = v[0]
	}

	pass, exists := u.User.Password()
	if exists {
		config["password"] = pass
	}

	return config
}
