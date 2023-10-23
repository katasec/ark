package scratchpad

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/katasec/ark/messaging"
)

var (
	m messaging.Messenger
)

func TestRedis(t *testing.T) {
	fmt.Println("Hello")
	connectionString := "redis://localhost:6379"
	queueName := "command-queue"

	m = messaging.NewRedisMessenger(connectionString, queueName)
	message := "Hello"

	// fmt.Printf("Press any key to send a message.\n")
	// fmt.Scanln()
	err := m.Send(queueName, message)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Message sent!")
	}

	// fmt.Printf("Press any key to receive message.\n")
	// fmt.Scanln()

	body, _, err := m.Receive()
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Received message!")
	}
	fmt.Println("Received message:" + body)
}

func TestRedisParseUrl(t *testing.T) {
	redisURL := "redis://localhost:6379"
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Addr:" + opt.Addr)
}
