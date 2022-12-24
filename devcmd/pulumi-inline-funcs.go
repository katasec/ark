package devcmd

import (
	"fmt"
	"strings"

	"github.com/katasec/ark/shell"
	pulumirunner "github.com/katasec/pulumi-runner"
	utils "github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createInlineProgram Creates an Inline Pulumi Program
func (d *DevCmd) createInlineProgram(pulumiFunc pulumi.RunFunc, stackName string) *pulumirunner.InlineProgram {

	logger := utils.ConfigureLogger(d.Config.LogFile)

	args := &pulumirunner.InlineProgramArgs{
		ProjectName: ProjectNamePrefix,
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

func getReference(stackFQDN string, key string) (output string, err error) {
	myCmd := fmt.Sprintf("pulumi stack -s %s output %s", stackFQDN, key)

	value, err := shell.ExecShellCmd(myCmd)
	value = strings.Trim(value, "\n")

	return value, err
}

func getDefaultPulumiOrg() (string, error) {

	value, err := shell.ExecShellCmd("pulumi org get-default")
	value = strings.Trim(value, "\n")

	return value, err
}
