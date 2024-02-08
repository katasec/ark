package prunner

import (
	"context"
	"log"
	"os"
	"strings"

	pulumirunner "github.com/katasec/pulumi-runner"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

// PRunner is a helper struct for running pulumi programs
type PRunner struct {
	ArkImage     string
	configdata   string
	resourceName string
	workDir      string
	localProgram pulumirunner.LocalProgram
}

func NewPRunner(arkImage string, configdata string, workDir string) *PRunner {
	runner := &PRunner{
		ArkImage:   arkImage,
		configdata: configdata,
		workDir:    workDir,
	}
	runner.setResourceName()
	runner.createLocalProgram()

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

func (p *PRunner) createLocalProgram() {

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
		WorkDir: p.workDir,
		Writer:  os.Stdout,
	}

	localProgram, err := pulumirunner.NewLocalProgram(args)
	if err != nil {
		log.Println(err.Error())
	} else {
		p.localProgram = *localProgram
	}

}

func (p *PRunner) Up() {
	p.localProgram.Stack.SetConfig(context.Background(), "arkdata", auto.ConfigValue{Value: p.configdata})
	p.localProgram.Up()
}

func (p *PRunner) Destroy() {
	p.localProgram.Stack.SetConfig(context.Background(), "arkdata", auto.ConfigValue{Value: p.configdata})
	p.localProgram.Destroy()
}
