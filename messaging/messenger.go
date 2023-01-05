package messaging

// Generic Messenger Interface
type Messenger interface {
	Send(message string) error
	Receive() (message string, subject string, err error)
	Complete() error
}
