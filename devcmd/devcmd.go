package devcmd

import (
	"fmt"
	"os"

	"github.com/dapr/cli/pkg/print"
	"github.com/katasec/ark/config"
)

type DevCmd struct {
	Config *config.Config
}

func NewDevCmd() *DevCmd {
	cfg := config.ReadConfig()

	return &DevCmd{
		Config: cfg,
	}
}

func (d *DevCmd) Setup() {
	message := fmt.Sprintf("Running Setup, please tail log file for more details %s", d.Config.LogFile)
	print.InfoStatusEvent(os.Stdout, message)
	fmt.Println()

	// Create Cloud Resources with Pulumi
	err := d.createLocal()

	// Update config file with links to new pulumi cloud resources.
	if err != nil {
		refreshConfig()
	}
}

func (d *DevCmd) Delete() {

	message := fmt.Sprintf("Running Delete, please tail log file for more details %s", d.Config.LogFile)
	print.InfoStatusEvent(os.Stdout, message)
	fmt.Println()

	// Delete Cloud Resources with Pulumi
	err := d.runWithProgressBar("Delete dev resources on Azure", addSbNsFunc, "dev", "destroy")

	// Update config
	if err != nil {
		refreshConfig()
	}

}
