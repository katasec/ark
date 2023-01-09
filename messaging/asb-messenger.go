package messaging

import (
	"context"
	"errors"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// AsbMessenger is an Azure Servicebus Implementation of the Messenger interface.
type AsbMessenger struct {
	ctx    context.Context
	client *azservicebus.Client

	queueName string

	sender          *azservicebus.Sender
	receiver        *azservicebus.Receiver
	ReceivedMessage *azservicebus.ReceivedMessage
}

// NewAsbMessenger Creates a new AsbMessenger which implements the Messenger interface
func NewAsbMessenger(connectionString string, queueName string) *AsbMessenger {

	// Create context for Azure Service Bus
	ctx := context.Background()

	client := getAsbClient(connectionString)

	// Create ASB Messenger struct
	asbMessenger := &AsbMessenger{
		ctx:    ctx,
		client: client,

		queueName: queueName,

		sender:   getAsbSender(client, queueName),
		receiver: getAsbReceiver(ctx, client, queueName),
	}

	return asbMessenger
}

func (m *AsbMessenger) Send(message string) error {

	// Send message
	err := m.sender.SendMessage(m.ctx, &azservicebus.Message{
		Body: []byte(message),
	}, nil)

	if err != nil {
		return err
	}

	return nil
}

func (m *AsbMessenger) Receive() (string, string, error) {

	message, messageBody, err := receiveMessage(m.ctx, m.receiver)
	if err != nil {
		return "", "", err
	}

	m.ReceivedMessage = message

	var subject string

	if message.Subject == nil {
		subject = ""
	} else {
		subject = *message.Subject
	}

	return messageBody, subject, err
}

func (m *AsbMessenger) Complete() error {

	err := m.receiver.CompleteMessage(m.ctx, m.ReceivedMessage, nil)

	if err != nil {
		var sbErr *azservicebus.Error

		if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeLockLost {
			// The message lock has expired. This isn't fatal for the client, but it does mean
			// that this message can be received by another Receiver (or potentially this one!).
			log.Printf("Message lock expired\n")

			// You can extend the message lock by calling receiver.RenewMessageLock(msg) before the
			// message lock has expired.
		}

		log.Println()
		return err
	}

	return nil
}

/*
-------------------------------------------------------
 Helper Funcs to interact with the Azure Service Bus
 via the Azure Service bus client
-------------------------------------------------------
*/

/*
getAsbClient creates an Azure Service Bus Client using the connection string
*/
func getAsbClient(connectionString string) *azservicebus.Client {

	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatal("Could not create servicebus Client: " + err.Error())
	}

	return client
}

/*
getAsbSender Generates a Service Bus Sender via the asb client
*/
func getAsbSender(client *azservicebus.Client, queueName string) *azservicebus.Sender {
	// Create sender
	sender, err := client.NewSender(queueName, nil)
	if err != nil {
		log.Fatal("Could not create sender !")
	}

	return sender
}

/*
getAsbReceiver Generates a Service Bus receiver via the asb client
*/
func getAsbReceiver(ctx context.Context, client *azservicebus.Client, queueName string) *azservicebus.Receiver {

	options := &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModePeekLock,
	}

	receiver, err := client.NewReceiverForQueue(queueName, options)
	if err != nil {
		log.Fatal("Could not create receiver: " + err.Error())
	}

	return receiver
}

/*
receiveMessage receives a message from the specified queue
*/
func receiveMessage(ctx context.Context, receiver *azservicebus.Receiver) (message *azservicebus.ReceivedMessage, messageBody string, err error) {

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	if err != nil {
		log.Println("Could not receive messages: " + err.Error())
		return nil, "", err
	}

	if messages != nil {
		if messages[0] == nil {
			log.Println("Could not receive messages: " + err.Error())
			return nil, "", err
		}
	} else {
		log.Printf("No message received\n")
		return nil, "", nil
	}

	message = messages[0]
	messageBody = string(message.Body)

	return message, messageBody, err
}
