package dev

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/katasec/ark/config"
	shell "github.com/katasec/utils/shell"
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

	var err error

	arkSpinner.InfoStatusEvent("Refershing config,please wait...")

	// Get config Details
	cfg := config.ReadConfig()

	// Set pulumi default org if not set
	if cfg.PulumiDefultOrg == "" {
		orgName, err := d.getDefaultPulumiOrg()
		utils.ExitOnError(err)
		cfg.PulumiDefultOrg = orgName
		cfg.Save()
	}

	// Contruct Pulumi Stack FQDN
	orgName := cfg.PulumiDefultOrg
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

	// Load cfg file and ensure any dirs are setup
	cfg := config.ReadConfig()
	cfg.SetupDirectories()

	// Use version in vars.go if unspecified in config file
	if cfg.DockerImages.Server == "" {
		cfg.DockerImages.Server = DEV_ARK_SERVER_IMAGE_NAME
		cfg.Save()
	}

	// Use version in vars.go if unspecified in config file
	if cfg.DockerImages.Worker == "" {
		cfg.DockerImages.Worker = DEV_ARK_WORKER_IMAGE_NAME
		cfg.Save()
	}

	var mounts []string
	arkSpinner.InfoStatusEvent("Starting Ark...")

	// Determine home folder path (for volume mounts)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// TODOS: Temporary hack, find alternative
	if runtime.GOOS != "windows" {
		shell.ExecShellCmd("chmod -R 644 " + path.Join(homeDir, ".pulumi", "credentials.json"))
	}

	if runtime.GOOS == "windows" {
		homeDir = winHomeDir(homeDir)
	}

	// ***************************************
	// Start Mysql Server
	// ***************************************
	envVars := []string{
		"MYSQL_ROOT_PASSWORD=Password123",
	}
	dh.StartContainerUI(DEV_ARK_DB_MYSQL_IMAGE_NAME, envVars, "3306", "db", nil, []string{}...)

	// ***************************************
	// Start Ark Server
	// ***************************************
	containerName := "arkserver"

	// using pipe as a separator for source and destination
	mounts = []string{
		fmt.Sprintf("%v|/root/.ark", config.GetArkDir()),
		// "/var/run/docker.sock|/root/.ark/",
		fmt.Sprintf("%v/.pulumi|/root/.pulumi", homeDir),
		fmt.Sprintf("%v/.azure|/root/.azure", homeDir),
	}

	envVars = []string{
		fmt.Sprintf("ASPNETCORE_URLS=http://%s:%s", cfg.ApiServer.Host, cfg.ApiServer.Port),
	}
	//md := []string{"/ark", "server"}
	dh.StartContainerUI(cfg.DockerImages.Server, envVars, cfg.ApiServer.Port, containerName, []string{"/ark", "server"}, mounts...)

	// ***************************************
	// Start Ark worker
	// ***************************************
	containerName = "arkworker"
	mounts = []string{
		fmt.Sprintf("%v/|/root/.ark", config.GetArkDir()),
		//fmt.Sprintf("%v/.ark|/home/app/.ark", homeDir),
		fmt.Sprintf("%v/.pulumi|/root/.pulumi", homeDir),
		//fmt.Sprintf("%v/.azure|/root/.azure", homeDir),
	}
	envVars = []string{
		"ARM_CLIENT_ID=" + os.Getenv("ARM_CLIENT_ID"),
		"ARM_CLIENT_SECRET=" + os.Getenv("ARM_CLIENT_SECRET"),
		"ARM_SUBSCRIPTION_ID=" + os.Getenv("ARM_SUBSCRIPTION_ID"),
		"ARM_LOCATION_NAME=" + os.Getenv("ARM_LOCATION_NAME"),
		"ARM_TENANT_ID=" + os.Getenv("ARM_TENANT_ID"),
	}
	dh.StartContainerUI(cfg.DockerImages.Worker, envVars, "0", containerName, []string{"/ark", "worker", "start"}, mounts...)

}

func (d *DevCmd) Stop() {
	// Load config file
	config := config.ReadConfig()

	// Stop Server
	dh.StopContainerUI(config.DockerImages.Server, "arkserver")

	// Stop worker
	dh.StopContainerUI(config.DockerImages.Worker, "arkworker")

}
func (d *DevCmd) Check() bool {
	return CheckSetupPreReqs()
}

func winHomeDir(homeDir string) string {
	// switch backslash to frontslash
	winDir := strings.Replace(homeDir, "\\", "/", -1)

	// remove colons
	winDir = strings.Replace(winDir, ":", "", -1)

	// And an "/mnt" in front
	winDir = path.Join("/mnt", winDir)
	fmt.Println(winDir)

	return winDir
}
