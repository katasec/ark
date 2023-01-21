package dev

import (
	"fmt"
	"log"
	"strings"

	"github.com/katasec/ark/shell"
	pulumirunner "github.com/katasec/pulumi-runner"
	utils "github.com/katasec/pulumi-runner/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createPulumiProgram Creates Pulumi Program using an in-line function passed as a parameter
func (d *DevCmd) createPulumiProgram(pulumiFn pulumi.RunFunc, stackName string) *pulumirunner.InlineProgram {

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
		PulumiFn: pulumiFn,
	}

	p, err := pulumirunner.NewInlineProgram(args)
	if err != nil {
		log.Fatal(err.Error())
	}

	return p
}

func (d *DevCmd) getReference(stackFQDN string, key string) (output string, err error) {
	myCmd := fmt.Sprintf("pulumi stack -s %s output %s", stackFQDN, key)

	value, err := shell.ExecShellCmd(myCmd)
	value = strings.Trim(value, "\n")

	return value, err
}

func (d *DevCmd) getDefaultPulumiOrg() (string, error) {

	value, err := shell.ExecShellCmd("pulumi org get-default")
	value = strings.Trim(value, "\n")

	return value, err
}
