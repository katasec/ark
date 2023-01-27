package worker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"encoding/json"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
	"github.com/katasec/ark/sdk/v0/messages"
	pulumirunner "github.com/katasec/pulumi-runner"
	"github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"gopkg.in/yaml.v2"
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

func jsonUnmarshall[V any](message string) (V, error) {
	var msg V
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Println("Invalid message:" + err.Error())
		log.Println("Bad message:" + message)
	}
	return msg, err
}

func yamlMarshall[V any](message V) (string, error) {
	// Convert message into yaml
	b, err := yaml.Marshal(message)
	if err != nil {
		fmt.Println("Could not covert request to yaml config data")
		log.Printf("Bad message: %v\n", message)
	}

	return string(b), err
}

// For e.g from subject  'createazurecloudspacerequest' or 'deleteazurecloudspacerequest', returns 'azurecloudspace'
func (w *Worker) getResourceName(subject string) string {
	resourceName := strings.ToLower(subject)
	resourceName = strings.Replace(resourceName, "delete", "", 1)
	resourceName = strings.Replace(resourceName, "create", "", 1)
	resourceName = strings.Replace(resourceName, "request", "", 1)
	return resourceName
}

// messageHandler Creates a pulumi program and injects the message as pulumi config
func (w *Worker) messageHandler(subject string, resourceName string, yamlconfig string) {

	// Create a pulumi program to handle this message
	p, err := w.createPulumiProgram(resourceName, messages.Runtimes.Dotnet)

	if err == nil {
		// Inject yaml config as input for pulumi program
		p.Stack.SetConfig(context.Background(), "arkdata", auto.ConfigValue{Value: yamlconfig})
		p.FixConfig()

		// Need code to check if another pulumi update is running
		// If yes then kill message and reject update.

		// Run pulumi up or destroy
		if strings.HasPrefix(strings.ToLower(subject), "delete") {
			p.Destroy()
		} else {
			p.Up()
		}
	} else {
		log.Println("Error creating pulumi program:" + err.Error())
	}

}

// createPulumiProgram creates a pulumi program from a git remote resource for the given resource name
func (w *Worker) createPulumiProgram(resourceName string, runtime string) (*pulumirunner.RemoteProgram, error) {

	logger := utils.ConfigureLogger(w.config.LogFile)
	projectPath := fmt.Sprintf("%s-handler", resourceName)

	log.Println("Project path:" + projectPath)
	args := &pulumirunner.RemoteProgramArgs{
		ProjectName: resourceName,
		GitURL:      "https://github.com/katasec/library.git",
		Branch:      "refs/remotes/origin/main",
		ProjectPath: projectPath,
		StackName:   "dev",
		Plugins: []map[string]string{
			{
				"name":    "azure-native",
				"version": "v1.89.1",
			},
		},
		Config: []map[string]string{
			{
				"name":  "azure-native:location",
				"value": "westus2",
			},
		},
		Runtime: runtime,
		Writer:  logger,
	}

	// Create a new pulumi program
	return pulumirunner.NewRemoteProgram(args)
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
		fmt.Println("**************************************")

		// Route the message
		subject = strings.ToLower(subject)
		resourceName := w.getResourceName(subject)

		switch resourceName {
		case "azurecloudspace":

			// Subject could be 'createazurecloudspacerequest' or 'deleteazurecloudspacerequest'

			// Convert json -> struct -> yaml and pass yaml as input to pulumi program
			msgStruct, _ := jsonUnmarshall[messages.AzureCloudspace](message)
			yamlConfig, _ := yamlMarshall(msgStruct)
			fmt.Println(yamlConfig)
			go w.messageHandler(subject, resourceName, yamlConfig)

		default:
			log.Printf("subject: %s", subject)
			log.Println("Unrecognized message, skipping")
		}

	}

}
