package prunner

import (
	"log"
	"os"
	"strings"

	pulumirunner "github.com/katasec/pulumi-runner"
)

// PRunner is a helper struct for running pulumi programs
type PRunner struct {
	ArkImage     string
	configdata   string
	resourceName string
}

func NewPRunner(arkImage string, configdata string) *PRunner {
	runner := &PRunner{
		ArkImage:   arkImage,
		configdata: configdata,
	}
	runner.setResourceName()

	return runner
}

func (p *PRunner) setResourceName() {

	// Image Name is after the last "/"
	words := strings.Split(p.ArkImage, "/")
	resourceName := words[len(words)-1]

	// Remove "version name"
	resourceName = strings.Split(resourceName, ":")[0]

	// Set Resource Name
	p.resourceName = resourceName
}

func (p *PRunner) Run() {

	log.Println("Do pulumi stuff!")

	args := &pulumirunner.LocalProgramArgs{
		ProjectName: p.resourceName,
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
		Writer: os.Stdout,
	}

	localProgram, err := pulumirunner.NewLocalProgram(args)
	if err != nil {
		log.Println(err.Error())
	} else {
		localProgram.Up()
	}

}

// // createPulumiProgram creates a pulumi program from a git remote resource for the given resource name
// func (p *PRunner) createPulumiProgram(resourceName string, runtime string) (*pulumirunner.RemoteProgram, error) {

// 	//logger := utils.ConfigureLogger(w.config.LogFile)
// 	//projectPath := fmt.Sprintf("%s-handler", resourceName)

// 	//log.Println("Project path:" + projectPath)

// 	args := &pulumirunner.

// 	// Create a new pulumi program
// 	return pulumirunner.NewRemoteProgram(args)
// }
