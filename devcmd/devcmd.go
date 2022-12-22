package devcmd

import (
	"fmt"
	"os"

	"github.com/dapr/cli/pkg/print"
	"github.com/hpcloud/tail"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/utils"
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

	//checkBeforeCreate()

	// Print link to log file
	message := fmt.Sprintf("Running Setup, please tail log file for more details %s", d.Config.LogFile)
	print.InfoStatusEvent(os.Stdout, message)

	// Start spinner
	note := "Seting up Azure components for dev environment"
	spinner := utils.NewArkSpinner()

	// Run pulumi up to create cloud resources
	spinner.Start(note)
	pulumi := d.createInlineProgram(setupAzureDeps, "dev")
	err := pulumi.Up()
	spinner.Stop(err, note)

	// Update config file with links to new pulumi cloud resources.
	if err != nil {
		refreshConfig()
	}
}

func (d *DevCmd) Delete() {

	// Print link to log file
	message := fmt.Sprintf("Running Delete, please tail log file for more details %s", d.Config.LogFile)
	print.InfoStatusEvent(os.Stdout, message)

	// Start spinner
	note := "Deleting Azure components from dev environment"
	spinner := utils.NewArkSpinner()

	// Run pulumi destroy to delete cloud resources
	spinner.Start(note)
	pulumi := d.createInlineProgram(setupAzureDeps, "dev")
	err := pulumi.Destroy()
	spinner.Stop(err, note)

	// Update config
	if err != nil {
		refreshConfig()
	}

}

func (d *DevCmd) Logs() {
	t, err := tail.TailFile(d.Config.LogFile, tail.Config{Follow: true})
	utils.ExitOnError(err)
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

func checkBeforeCreate() {
	fmt.Println("Check for existing ")
}
