package worker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"encoding/json"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
	"github.com/katasec/ark/sdk/v0/resources"
	pulumirunner "github.com/katasec/pulumi-runner"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
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

		// Route the message
		subject = strings.ToLower(subject)

		switch subject {
		case "azurecloudspace":
			go runAzureCloudspace(subject, message)
			//w.mq.Complete()
		default:
			log.Printf("subject: %s", subject)
			log.Println("Unrecognized message, skipping")
			//w.mq.Complete()
		}

	}
}

func runAzureCloudspace(subject string, message string) {

	// Convert request to struct
	msg := resources.AzureCloudspace{}
	json.Unmarshal([]byte(message), &msg)

	// Output for debug purposes
	log.Printf("Hub Name:" + msg.Hub.Name)

	// Create a pulumi program to handle this message
	p := createPulumiProgram(subject, resources.Runtimes.Dotnet)

	// Inject message details as input for pulumi program
	ctx := context.Background()
	p.Stack.SetConfig(ctx, "arkdata", auto.ConfigValue{Value: string(message)})

	// Need code to check if another pulumi update is running
	// If yes then kill message and reject update.

	//p.Stack.
	p.Up()

}
func createPulumiProgram(subject string, runtime string) *pulumirunner.RemoteProgram {

	projectPath := fmt.Sprintf("%s-handler", subject)

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
	}

	// Create a new pulumi program
	return pulumirunner.NewRemoteProgram(args)
}
