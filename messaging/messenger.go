package messaging

// Generic Messenger Interface
type Messenger interface {
	Send(message string) error
	Receive() (string, error)
	Complete() error
}
