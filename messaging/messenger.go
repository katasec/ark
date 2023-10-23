package messaging

// Generic Messenger Interface
type Messenger interface {
	Send(subject string, message string) error
	Receive() (message string, subject string, err error)
	//Complete() error
}
