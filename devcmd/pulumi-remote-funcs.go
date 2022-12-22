package devcmd

import (
	"os"
	"path/filepath"

	"github.com/dapr/cli/pkg/print"
	pulumirunner "github.com/katasec/pulumi-runner"
	utils "github.com/katasec/pulumi-runner/utils"
)

func initRemoteProgram() *pulumirunner.RemoteProgram {

	homedir, _ := os.UserHomeDir()
	logger := utils.ConfigureLogger(filepath.Join(homedir + "ark.log"))

	args := &pulumirunner.RemoteProgramArgs{
		ProjectName: "ark-init",
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

func createRemote() {

	message := "Creating cloud resources for local dev environment."
	stopSpinning := print.Spinner(os.Stdout, message)
	p := initRemoteProgram()
	p.Up()

	stopSpinning(print.Success)
	//print.FailureStatusEvent(os.Stdout, message)
}

func deleteRemote() {
	p := initRemoteProgram()

	message := "Deleting cloud resources for local dev environment."
	stopSpinning := print.Spinner(os.Stdout, message)

	p.Destroy()
	stopSpinning(print.Success)

}
