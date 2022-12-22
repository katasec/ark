package devcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dapr/cli/pkg/print"
	pulumirunner "github.com/katasec/pulumi-runner"
	utils "github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createInlineProgram Creates an Inline Pulumi Program
func (d *DevCmd) createInlineProgram(pulumiFunc pulumi.RunFunc, stackName string) *pulumirunner.InlineProgram {

	logger := utils.ConfigureLogger(d.Config.LogFile)

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

func (d *DevCmd) runWithProgressBar(message string, pulumiFunc pulumi.RunFunc, stackName string, action string) error {

	var err error

	// Show status
	//stopSpinning := print.Spinner(os.Stdout, message)

	// Initializer an in-line pulumi program with the passed in func
	p := d.createInlineProgram(pulumiFunc, stackName)

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
	// if err != nil {
	// 	stopSpinning(print.Failure)
	// } else {
	// 	stopSpinning(print.Success)
	// }

	return err
}

func (d *DevCmd) createLocal() error {
	runSuccess := true

	// // Setup Ark Resource Group
	// err := d.runWithProgressBar("Setup Azure resource group", createRgFunc, "dev", "up")
	// if err != nil {
	// 	runSuccess = false
	// }
	// // Setup Ark  Storage Account
	// err = d.runWithProgressBar("Setup Azure storage account", addStrgFunc, "dev", "up")
	// if err != nil {
	// 	runSuccess = false
	// }
	// // Setup Ark Service Bus Name Space
	// err = d.runWithProgressBar("Setup Azure Service Bus Namespace", addSbNsFunc, "dev", "up")
	// if err != nil {
	// 	runSuccess = false
	// }

	// Setup Ark Service Bus Name Space
	err := d.runWithProgressBar("Setup Azure command queue", setupAzureDeps, "dev", "up")
	if err != nil {
		runSuccess = false
	}

	if !runSuccess {
		fmt.Println()
		print.InfoStatusEvent(os.Stdout, "One or more of the above had errors. Please check ark logs.")
	}

	return err
}
