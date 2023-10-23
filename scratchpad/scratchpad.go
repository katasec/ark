package scratchpad

import (
	"fmt"
	"log"

	"github.com/katasec/ark/messaging"
)

var (
	m messaging.Messenger
)

func Start() {
	testRedisMessenger()
}

func testRedisMessenger() {
	connectionString := "redis://localhost:6379"
	queueName := "command-queue"

	m := messaging.NewRedisMessenger(connectionString, queueName)
	message := "Hello"

	fmt.Printf("Press any key to send a message.\n")
	fmt.Scanln()
	err := m.Send(queueName, message)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Message sent!")
	}

	fmt.Printf("Press any key to receive message.\n")
	fmt.Scanln()
	body, _, err := m.Receive()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Received message:" + body)
}
