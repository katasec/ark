package messaging

// Generic Messenger Interface
type Messenger interface {
	Send(queueName string, message string) error
	Receive(queueName string) (string, error)
	Complete() error
}
