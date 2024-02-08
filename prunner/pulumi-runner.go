package prunner

import "log"

// import (
// pulumirunner "github.com/katasec/pulumi-runner"
// )

// PRunner is a helper struct for running pulumi programs
type PRunner struct {
	ArkImage   string
	configdata string
}

func NewPRunner(arkImage string, configdata string) *PRunner {
	runner := &PRunner{
		ArkImage:   arkImage,
		configdata: configdata,
	}

	return runner
}

func (p *PRunner) Run() {
	log.Println("Do pulumi stuff!")
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
