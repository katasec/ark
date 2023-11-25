package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	resources "github.com/katasec/ark/resources"
	shell "github.com/katasec/utils/shell"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
	pulumirunner "github.com/katasec/pulumi-runner"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type Worker struct {
	config *config.Config
	respQ  messaging.Messenger
	cmdQ   messaging.Messenger
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
	w.setupAzureCreds()
	log.Println("Starting worker")

	// Inifinite loop polling messages
	for {
		// This is a blocking receive
		log.Println("polling for message...")
		message, subject, err := w.cmdQ.Receive()
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
		executeCommand(resourceName, w, subject, message, err)
	}
}

// messageHandler Creates a pulumi program and injects the message as pulumi config
func (w *Worker) messageHandler(subject string, resourceName string, yamlconfig string, c chan error) {

	log.Println("Before creating pulumi program")

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

	//logger := utils.ConfigureLogger(w.config.LogFile)
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

func (w *Worker) setupAzureCreds() {

	log.Println("Reading Azure creds from env vars")
	envvars := []string{"ARM_CLIENT_ID", "ARM_CLIENT_SECRET", "ARM_TENANT_ID", "ARM_SUBSCRIPTION_ID"}

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

	if !ok {
		log.Println("Some env vars not set, exitting...")
		os.Exit(1)
	}
}

func executeCommand(resourceName string, w *Worker, subject string, message string, err error) {
	switch resourceName {
	case "azurecloudspace":

		c := make(chan error)
		go w.messageHandler(subject, resourceName, message, c)
		handlerError := <-c

		if handlerError == nil {
			fmt.Println("Handler ran without errors !")

			w.respQ.Send(subject, message)

			if strings.HasPrefix(strings.ToLower(subject), "delete") {
				fmt.Println("TODO: Delete Cloudpace from DB")
			} else {
				fmt.Println("TODO: Add Cloudpace to DB")
			}
		} else {
			fmt.Println("Handler errors:" + handlerError.Error())
		}

	case "hellosuccess":

		c := make(chan error)
		go w.messageHandler(subject, resourceName, message, c)

	default:
		log.Printf("subject: %s", subject)
		log.Println("Unrecognized message, skipping")
	}
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
