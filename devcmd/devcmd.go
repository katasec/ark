package devcmd

import (
	"fmt"
	"os"

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
	fmt.Println("Checking pre-requisites before running setup:")
	if !CheckSetupPreReqs() {
		os.Exit(1)
	}

	// Setup spinner
	spinner := utils.NewArkSpinner()

	// Print link to log file
	message := fmt.Sprintf("Running Setup, please tail log file for more details %s", d.Config.LogFile)
	spinner.InfoStatusEvent(message)

	// Run pulumi up to create cloud resources
	note := "Seting up Azure components for dev environment"
	spinner.Start(note)
	pulumi := d.createInlineProgram(setupAzureComponents, StackName)
	err := pulumi.Up()
	spinner.Stop(err, note)

	// Update config file with links to new pulumi cloud resources.
	if err == nil {
		refreshConfig()
	}

}

func (d *DevCmd) Delete() {

	// Setup spinner
	spinner := utils.NewArkSpinner()

	// Print link to log file
	message := fmt.Sprintf("Running Delete, please tail log file for more details %s", d.Config.LogFile)
	spinner.InfoStatusEvent(message)

	// Run pulumi destroy to delete cloud resources
	note := "Deleting Azure components from dev environment"
	spinner.Start(note)
	pulumi := d.createInlineProgram(setupAzureComponents, "dev")
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
