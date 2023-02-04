package dev

import (
	"fmt"
	"log"
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

	// Load config file and ensure any dirs are setup
	config := config.ReadConfig()
	config.SetupDirectories()

	// Use version in vars.go if unspecified in config file
	if config.DockerImages.Server == "" {
		config.DockerImages.Server = DEV_ARK_SERVER_IMAGE_NAME
		config.Save()
	}

	// Use version in vars.go if unspecified in config file
	if config.DockerImages.Worker == "" {
		config.DockerImages.Worker = DEV_ARK_WORKER_IMAGE_NAME
		config.Save()
	}

	var mounts []string
	arkSpinner.InfoStatusEvent("Starting Ark...")

	// Determine home folder path (for volume mounts)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// ***************************************
	// Start Ark Server
	// ***************************************
	containerName := "arkserver"
	mounts = []string{
		fmt.Sprintf("%v/.ark:/home/app/.ark", homeDir),
	}
	envVars := []string{
		fmt.Sprintf("ASPNETCORE_URLS=http://%s:%s", config.ApiServer.Host, config.ApiServer.Port),
	}
	dh.StartContainerUI(config.DockerImages.Server, envVars, config.ApiServer.Port, containerName, nil, mounts...)

	// ***************************************
	// Start Ark worker
	// ***************************************
	containerName = "arkworker"
	mounts = []string{
		fmt.Sprintf("%v/.ark:/root/.ark", homeDir),
		fmt.Sprintf("%v/.pulumi:/root/.pulumi", homeDir),
		fmt.Sprintf("%v/.azure:/root/.azure", homeDir),
	}
	dh.StartContainerUI(config.DockerImages.Worker, nil, "0", containerName, []string{"/ark", "worker", "start"}, mounts...)

	// envVars := []string{
	// 	"POSTGRES_USER=" + DevDbDefaultUser,
	// 	"POSTGRES_PASSWORD=" + DevDbDefaultPass,
	// }
	// port := "5432"
}

func (d *DevCmd) Check() bool {
	return CheckSetupPreReqs()
}
