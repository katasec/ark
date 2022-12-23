package devcmd

import (
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
