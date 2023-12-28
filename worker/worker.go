package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/katasec/ark/requests"
	resources "github.com/katasec/ark/resources"
	shell "github.com/katasec/utils/shell"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
	pulumirunner "github.com/katasec/pulumi-runner"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"

	logx "github.com/katasec/ark/log"
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
		logx.Logger.Info("The subject was:" + subject)

		// Execute command based on subject
		switch subject {
		case "createazurecloudspacerequest":
			executeCommand[requests.CreateAzureCloudspaceRequest](w, message, err)
		case "deleteazurecloudspacerequest":
			executeCommand[requests.DeleteAzureCloudspaceRequest](w, message, err)
		}

	}
}

func executeCommand[T requests.RequestInterface](w *Worker, payload string, err error) error {

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
	go w.messageHandler(action, resourceName, payload, c)
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

// messageHandler Creates a pulumi program and injects the message as pulumi config
func (w *Worker) messageHandler(action string, resourceName string, yamlconfig string, c chan error) {

	log.Println("Before creating pulumi program")

	// Create a pulumi program to handle this message
	p, err := w.createPulumiProgram(resourceName, resources.Runtimes.Dotnet)

	if err == nil {
		// Inject yaml config as input for pulumi program
		p.Stack.SetConfig(context.Background(), "arkdata", auto.ConfigValue{Value: yamlconfig})
		p.FixConfig()

		// Run pulumi up or destroy
		if strings.HasPrefix(strings.ToLower(action), "delete") {
			_, err := p.Destroy()
			c <- err
		} else {
			_, err := p.Up()
			c <- err
		}
	} else {
		log.Println("Error creating pulumi program:" + err.Error())
		c <- err
	}
}

// createPulumiProgram creates a pulumi program from a git remote resource for the given resource name
func (w *Worker) createPulumiProgram(resourceName string, runtime string) (*pulumirunner.RemoteProgram, error) {

	//logger := utils.ConfigureLogger(w.config.LogFile)
	//projectPath := fmt.Sprintf("%s-handler", resourceName)

	//log.Println("Project path:" + projectPath)

	args := &pulumirunner.RemoteProgramArgs{
		ProjectName: resourceName,
		ProjectPath: resourceName + "-handler", //projectPath,
		GitURL:      "https://github.com/katasec/library.git",
		Branch:      "refs/remotes/origin/main",
		StackName:   "dev",
		Plugins: []map[string]string{
			{
				"name":    "azure-native",
				"version": "v2.7.0",
			},
		},
		Config: []map[string]string{
			{
				"name":  "azure-native:location",
				"value": "westus2",
			},
		},
		Runtime: runtime,
		Writer:  os.Stdout,
	}

	// Create a new pulumi program
	return pulumirunner.NewRemoteProgram(args)
}

// getAzureCredsFromEnv reads Azure creds from env vars and sets them in pulumi config
func (w *Worker) getAzureCredsFromEnv() {

	// Define env vars to read
	log.Println("Reading Azure creds from env vars")
	envvars := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_TENANT_ID",
		"ARM_SUBSCRIPTION_ID",
	}

	// Check if all env vars are set
	ok := false
	for _, envvar := range envvars {
		if os.Getenv(envvar) != "" {
			ok = true
			shell.ExecShellCmd("pulumi config set azure-native:clientId " + os.Getenv(envvar))
		} else {
			ok = false
			log.Println("Env var " + envvar + " not set")
		}
	}

	// Exit if any env var is not set
	if !ok {
		log.Println("Some env vars not set, exitting...")
		os.Exit(1)
	}
}
