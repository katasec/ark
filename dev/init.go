package dev

import (
	"os"

	"github.com/dapr/cli/pkg/print"
	// pulumihelper "github.com/katasec/pulumi-helper"

	pulumirunner "github.com/katasec/pulumi-runner"
	utils "github.com/katasec/pulumi-runner/utils"
)

func createPulumiProgram() *pulumirunner.RemoteProgram {

	logger := utils.ConfigureLogger("ark.log")

	args := &pulumirunner.RemoteProgramArgs{
		ProjectName: "ArkInit",
		GitURL:      "https://github.com/katasec/ArkInit.git",
		Branch:      "refs/remotes/origin/main",
		ProjectPath: "Azure",
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
		Runtime: "dotnet",
		Writer:  logger,
	}

	return pulumirunner.NewRemoteProgram(args)
}
func Create() {

	message := "Creating cloud resources for local dev environment."
	stopSpinning := print.Spinner(os.Stdout, message)
	p := createPulumiProgram()
	p.Up()

	stopSpinning(print.Success)
	//print.FailureStatusEvent(os.Stdout, message)
}

func Delete() {
	p := createPulumiProgram()

	message := "Deleting cloud resources for local dev environment."
	stopSpinning := print.Spinner(os.Stdout, message)

	p.Destroy()
	stopSpinning(print.Success)

}
