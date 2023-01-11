package worker

import (
	"fmt"
	"log"
	"strings"

	"encoding/json"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
)

type Worker struct {
	config *config.Config
	mq     messaging.Messenger
}

func NewWorker() *Worker {

	// Read from local config  file
	cfg := config.ReadConfig()

	// Get queue name and access creds from config
	connectionString := cfg.AzureConfig.MqConfig.MqConnectionString
	queueName := cfg.AzureConfig.MqConfig.MqName

	// Create an mq client
	var mq messaging.Messenger = messaging.NewAsbMessenger(connectionString, queueName)

	fmt.Println("queueName is:" + queueName)
	// Return worker with local config
	return &Worker{
		config: cfg,
		mq:     mq,
	}
}

func (w *Worker) Start() {
	log.Println("Starting worker")

	// Inifinite loop polling messages
	for {
		// This is a blocking receive
		log.Println("Waiting for message...")
		message, subject, err := w.mq.Receive()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// Log Message
		log.Println("The subject was:" + subject)

		log.Println("The message was:" + message)

		// Route the message
		log.Println("Routing message to handler")
		switch strings.ToLower(subject) {
		case "azurecloudspacerequest":
			go w.CloudSpaceRequestHandler(message)
		default:
			log.Println("Unrecognised subject: '" + subject + "', completing message.")
			w.mq.Complete()
		}

		//go w.processMessage(req)

	}
}

func (w *Worker) CloudSpaceRequestHandler(message string) {

	log.Println("CloudSpaceRequestHandler handling message")

	c := &CloudSpaceRequest{}
	json.Unmarshal([]byte(message), c)

	log.Printf("The project name is:" + c.ProjectName)
	log.Printf("The date  is:" + c.DtTimeStamp.Format("2006-01-02 15:04:05"))

	log.Printf("Completing message")
	w.mq.Complete()
}
