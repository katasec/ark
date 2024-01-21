package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/katasec/ark/requests"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
)

type Worker struct {
	config *config.Config

	// respQ is the queue to send responses to the server
	respQ messaging.Messenger

	// cmdQ is the queue to receive commands from the server
	cmdQ messaging.Messenger
}

func NewWorker() *Worker {

	// Read from local config  file
	cfg := config.ReadConfig()

	// Return worker with local config
	return &Worker{
		config: cfg,
		cmdQ:   messaging.NewRabbitMqMessenger(cfg.MqConnStr, cfg.CmdQ),
		respQ:  messaging.NewRabbitMqMessenger(cfg.MqConnStr, cfg.RespQ),
	}
}

func (w *Worker) Start() {

	// Read Azure creds on startup
	w.getAzureCredsFromEnv()
	log.Println("Starting worker")

	// Inifinite loop polling messages
	for {
		// This is a blocking receive waiting for messages from the server
		log.Println("polling for message...")
		message, subject, err := w.cmdQ.Receive()
		if err != nil {
			log.Println("Infinite loop polling for message, error:" + err.Error())
			continue
		}

		// Log Message
		subject = strings.ToLower(subject)
		log.Println("Worker received subject was:" + subject)

		// Execute command based on subject
		switch subject {
		case "createazurecloudspacerequest":
			executeCommand[requests.CreateAzureCloudspaceRequest](w, message, err)
		case "deleteazurecloudspacerequest":
			executeCommand[requests.DeleteAzureCloudspaceRequest](w, message, err)
		case "createhellorequest":
			executeCommand[requests.CreateHelloRequest](w, message, err)
		default:
			log.Println("Invalid subject:" + subject)
		}

	}
}

func executeCommand[T requests.RequestInterface](w *Worker, payload string, err error) error {

	// Output command name
	var x T
	requestName := reflect.TypeOf(x).Name()
	fmt.Println("Executing command:" + requestName)

	// Convert payload to message type
	var message T
	json.Unmarshal([]byte(payload), &message)
	if err != nil {
		log.Println("Error unmarshalling message:" + err.Error())
		return err
	}

	// Extract action and resource name from request
	action := message.GetActionType()
	resourceName := message.GetResourceType()

	// Run handler in a go routine to create/destroy infra
	c := make(chan error)

	// Run Pulumi handler
	//go w.pulumiHandler(action, resourceName, payload, c)

	// Run Terraform handler
	go w.terraformHandler(action, resourceName, payload, c)

	handlerError := <-c

	// Send result to server via response queue on success
	if handlerError == nil {
		fmt.Println("Handler ran without errors !")
		fmt.Println("Sending response to server:" + requestName)
		w.respQ.Send(requestName, payload)
	} else {
		fmt.Println("Handler errors:" + handlerError.Error())
		return handlerError
	}

	return nil
}
