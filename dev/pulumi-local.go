package dev

import (
	"fmt"
	"os"
	"strings"

	"github.com/dapr/cli/pkg/print"
	pulumirunner "github.com/katasec/pulumi-runner"
	utils "github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func initLocalProgram(pulumiFunc pulumi.RunFunc, stackName string) *pulumirunner.InlineProgram {

	homedir, _ := os.UserHomeDir()
	logger := utils.ConfigureLogger(homedir + "/ark.log")

	args := &pulumirunner.InlineProgramArgs{
		ProjectName: "ark-init",
		StackName:   stackName,
		Config: []map[string]string{
			{
				"name":  "azure-native:location",
				"value": "westus2",
			},
		},
		Writer:   logger,
		PulumiFn: pulumiFunc,
	}

	return pulumirunner.NewInlineProgram(args)
}

func createLocal() {

	runSuccess := true
	var err error

	// Setup Ark Resource Group
	err = runWithProgressBar("Setup Ark resource group", createResourceGroup, "resource-group", "up")
	if err != nil {
		runSuccess = false
	}
	// Setup Ark  Storage Account
	err = runWithProgressBar("Setup Ark storage account", createStorageAccount, "storage-account", "up")
	if err != nil {
		runSuccess = false
	}
	// Setup Ark Service Bus Name Space
	err = runWithProgressBar("Setup Ark Service Bus Name Space", createSbNameSpace, "service-bus", "up")
	if err != nil {
		runSuccess = false
	}

	if !runSuccess {
		fmt.Println()
		print.InfoStatusEvent(os.Stdout, "One or more of the above had errors. Please check ark logs.")
	}
}

func deleteLocal() {

	// Destroy Resource Group
	runWithProgressBar("Setup Ark resource group", createResourceGroup, "resource-group", "destroy")

	// Destroy Storage Account
	runWithProgressBar("Setup Ark storage account", createStorageAccount, "storage-account", "destroy")

	// Destroy Ark Service Bus Name Space
	runWithProgressBar("Setup Ark Service Bus Name Space", createSbNameSpace, "service-bus", "destroy")

}

func runWithProgressBar(message string, pulumiFunc pulumi.RunFunc, stackName string, action string) error {

	var err error

	// Show status
	stopSpinning := print.Spinner(os.Stdout, message)

	// Initializer an in-line pulumi program with the passed in func
	p := initLocalProgram(pulumiFunc, stackName)

	// Call the appropriate action on the program
	action = strings.ToLower(action)
	switch action {
	case "preview":
		err = p.Preview()
	case "up":
		err = p.Up()
	case "destroy":
		err = p.Destroy()
	}

	// Output graphic status
	if err != nil {
		stopSpinning(print.Failure)
	} else {
		stopSpinning(print.Success)
	}

	return err
}
