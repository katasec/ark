package worker

import (
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
		fmt.Println("**************************************")
		log.Println("The subject was:" + subject)
		fmt.Println("**************************************")

		// Route the message
		subject = strings.ToLower(subject)

		switch subject {
		case "createazurecloudspacerequest":
			go w.AzureCloudspaceHandler(subject, message)
		case "deleteazurecloudspacerequest":
			go w.AzureCloudspaceHandler(subject, message)
		default:
			log.Printf("subject: %s", subject)
			log.Println("Unrecognized message, skipping")
		}

	}

}

func (w *Worker) AzureCloudspaceHandler(subject string, message string) {

	// Convert request to struct
	msg := messages.AzureCloudspace{}
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Output for debug purposes
	log.Printf("Hub Name:" + msg.Hub.Name)

	// Convert message into yaml
	yamlconfig, err := yaml.Marshal(msg)
	if err != nil {
		fmt.Println("Could not covert request to yaml config data")
	}

	// Create a pulumi program to handle this message
	resourceName := "azurecloudspace"
	p, err := w.createPulumiProgram(resourceName, messages.Runtimes.Dotnet)

	if err != nil {
		// Inject message details as input for pulumi program
		//ctx := context.Background()
		//p.Stack.SetConfig(ctx, "arkdata", auto.ConfigValue{Value: string(message)})
		//p.Stack.SetConfig(ctx, "arkdata", auto.ConfigValue{Value: string(yamlconfig)})
		p.SetConfig(auto.ConfigValue{Value: string(yamlconfig)})
		p.FixConfig()
		// Need code to check if another pulumi update is running
		// If yes then kill message and reject update.

		// Run pulumi up or destroy
		if subject == "deleteazurecloudspacerequest" {
			p.Destroy()
		} else {
			p.Up()
		}
	} else {
		log.Printf("Error creating pulumi program: %+v\n", err.Error())
	}

}

func (w *Worker) createPulumiProgram(resourceName string, runtime string) (*pulumirunner.RemoteProgram, error) {

	logger := utils.ConfigureLogger(w.config.LogFile)
	projectPath := fmt.Sprintf("%s-handler", resourceName)

	log.Println("Project path:" + projectPath)
	args := &pulumirunner.RemoteProgramArgs{
		ProjectName: "katasec-go-helloworld",
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
