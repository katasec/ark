package worker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/katasec/ark/client"
	resources "github.com/katasec/ark/resources/v0"

	"encoding/json"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
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

// Unmarshals a string to the provided type 'V'
func jsonUnmarshall[T any](message string) (T, error) {
	//log.Println("The message was:" + message)
	var msg T
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Println("Invalid message:" + err.Error())
		log.Println("jsonUnmarshall error:" + message)
	}
	return msg, err
}

// Marshals a struct of type 'V' to a yaml string
func yamlMarshall[T any](message T) (string, error) {
	// Convert message into yaml
	b, err := yaml.Marshal(message)
	if err != nil {
		fmt.Println("Could not covert request to yaml config data")
		log.Printf("yamlMarshall error: %v\n", message)
	}

	return string(b), err
}

func jsonToYaml[T any](message string) (string, error) {

	myStruct, _ := jsonUnmarshall[T](message)
	yamlString, err := yamlMarshall(myStruct)
	if err != nil {
		fmt.Printf("There were errors, this message will not be processed. Message: %s\n", message)
	}

	return yamlString, err
}

// For e.g from subject  'createazurecloudspacerequest' or 'deleteazurecloudspacerequest', returns 'azurecloudspace'
func (w *Worker) getResourceName(subject string) string {
	fmt.Println("getResourceName():" + subject)
	resourceName := strings.ToLower(subject)
	resourceName = strings.Replace(resourceName, "delete", "", 1)
	resourceName = strings.Replace(resourceName, "create", "", 1)
	resourceName = strings.Replace(resourceName, "request", "", 1)
	return resourceName
}

type autoResult struct {
	result any
	err    error
}

// messageHandler Creates a pulumi program and injects the message as pulumi config
func (w *Worker) messageHandler(subject string, resourceName string, yamlconfig string, c chan error) {

	// Create a pulumi program to handle this message
	p, err := w.createPulumiProgram(resourceName, resources.Runtimes.Dotnet)

	if err == nil {
		// Inject yaml config as input for pulumi program
		p.Stack.SetConfig(context.Background(), "arkdata", auto.ConfigValue{Value: yamlconfig})
		p.FixConfig()

		// Run pulumi up or destroy
		if strings.HasPrefix(strings.ToLower(subject), "delete") {
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
			log.Println("Infinite loop polling for message, error:" + err.Error())
			continue
		}

		subject = strings.ToLower(subject)
		fmt.Println("Received Subject:" + subject)

		// Log Message
		log.Println("The subject was:" + subject)

		// Route the message by resource name
		resourceName := w.getResourceName(subject)
		switch resourceName {
		case "azurecloudspace":
			// Convert json -> struct -> yaml and pass yaml as input to pulumi program
			if yamlConfig, err := jsonToYaml[resources.AzureCloudspace](message); err == nil {
				c := make(chan error)
				go w.messageHandler(subject, resourceName, yamlConfig, c)
				handlerError := <-c

				//  If Handler ran succesfully, update DB
				if handlerError == nil {

					fmt.Println("Handler ran without errors !")
					// Create ark api client
					arkClient := client.NewArkClient()

					// Convert nmessage to Azurecloudspace
					cs, err := jsonUnmarshall[resources.AzureCloudspace](message)
					if err != nil {
						break
					}

					// Update DB
					if strings.HasPrefix(strings.ToLower(subject), "delete") {
						arkClient.DeleteCloudSpace(cs)
					} else {
						arkClient.AddCloudSpace(cs)
					}
				} else {
					fmt.Println("Handler errors:" + handlerError.Error())
				}

			}
		case "hellosuccess":
			// Convert json -> struct -> yaml and pass yaml as input to pulumi program
			if yamlConfig, err := jsonToYaml[resources.HelloSuccess](message); err == nil {
				c := make(chan error)
				go w.messageHandler(subject, resourceName, yamlConfig, c)
			}

		default:
			log.Printf("subject: %s", subject)
			log.Println("Unrecognized message, skipping")
		}

	}

}
