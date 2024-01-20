package worker

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/katasec/ark/resources"
	pulumirunner "github.com/katasec/pulumi-runner"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

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

// pulumiHandler Creates a pulumi program and injects the message as pulumi config
func (w *Worker) pulumiHandler(action string, resourceName string, configdata string, c chan error) {

	log.Println("Before creating pulumi program")

	// Create a pulumi program to handle this message
	p, err := w.createPulumiProgram(resourceName, resources.Runtimes.Dotnet)
	if err != nil {
		log.Println("Error creating pulumi program:" + err.Error())
		c <- err
		return
	}

	// Inject yaml config as input for pulumi program
	p.Stack.SetConfig(context.Background(), "arkdata", auto.ConfigValue{Value: configdata})
	p.FixConfig()

	// Run pulumi up or destroy
	if strings.HasPrefix(strings.ToLower(action), "delete") {
		_, err := p.Destroy()
		c <- err
	} else {
		_, err := p.Up()
		c <- err
	}
}
