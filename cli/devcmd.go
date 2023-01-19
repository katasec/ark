package cli

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
	pulumi := d.createPulumiProgram(setupAzureComponents, StackName)
	err := pulumi.Up()
	spinner.Stop(err, note)

	// Update config file with links to new pulumi cloud resources.
	if err == nil {
		d.RefreshConfig()
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
	pulumi := d.createPulumiProgram(setupAzureComponents, "dev")
	err := pulumi.Destroy()
	spinner.Stop(err, note)

	// Update config
	if err != nil {
		d.RefreshConfig()
	}

}

func (d *DevCmd) Logs() {
	t, err := tail.TailFile(d.Config.LogFile, tail.Config{Follow: true})
	utils.ExitOnError(err)
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

func (d *DevCmd) RefreshConfig() {

	arkSpinner.InfoStatusEvent("Refershing config,please wait...")

	// Get config Details
	cfg := config.ReadConfig()

	// Contruct Pulumi Stack FQDN
	orgName, err := d.getDefaultPulumiOrg()
	utils.ExitOnError(err)
	stackFQDN := fmt.Sprintf("%s/%s/%s", orgName, ProjectNamePrefix, StackName)

	/*
	 Extract Azure resource details from Pulumi Exports
	*/

	// Resource Group Name
	cfg.AzureConfig.ResourceGroupName, err = d.getReference(stackFQDN, ResourceGroupName)
	utils.ExitOnError(err)

	// Mq Connection String
	cfg.AzureConfig.MqConfig.MqConnectionString, err = d.getReference(stackFQDN, MqConnectionString)
	utils.ExitOnError(err)

	// Mq Connection String
	cfg.AzureConfig.MqConfig.MqName, err = d.getReference(stackFQDN, CommandQueueName)
	utils.ExitOnError(err)

	// Log Storage Account Name
	cfg.AzureConfig.StorageConfig.LogStorageAccountName, err = d.getReference(stackFQDN, LogStorageAccountName)
	utils.ExitOnError(err)

	// Log Storage Account Endpoint
	cfg.AzureConfig.StorageConfig.LogStorageEndpoint, err = d.getReference(stackFQDN, LogStorageEndpoint)
	utils.ExitOnError(err)

	// Log Storage LogStorageKey
	cfg.AzureConfig.StorageConfig.LogStorageKey, err = d.getReference(stackFQDN, LogStorageKey)
	utils.ExitOnError(err)

	// Log Storage Container
	cfg.AzureConfig.StorageConfig.LogsContainer, err = d.getReference(stackFQDN, LogContainerName)
	utils.ExitOnError(err)

	// Pulumi State Container
	cfg.AzureConfig.StorageConfig.PulumiStateContainer, err = d.getReference(stackFQDN, PulumiStateContainerName)
	utils.ExitOnError(err)

	// Save Azure resource details to config file
	cfg.Save()

}

func (d *DevCmd) Start() {
	fmt.Println("Starting awesomeness!")

	// Start Ark Container
	imageName := DEV_ARK_IMAGE_NAME
	// envVars := []string{
	// 	"POSTGRES_USER=" + DevDbDefaultUser,
	// 	"POSTGRES_PASSWORD=" + DevDbDefaultPass,
	// }
	// port := "5432"

	dh.ContainerRemove("/arkserver")
	dh.StartContainerUI(imageName, nil, "", "arkserver", nil)
}

func (d *DevCmd) Check() bool {
	return CheckSetupPreReqs()
}
