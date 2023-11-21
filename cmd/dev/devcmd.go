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
	"github.com/katasec/ark/utils"
	shell "github.com/katasec/utils/shell"
)

type DevCmd struct {
	Config *config.Config
}

func NewDevCmd() *DevCmd {
	//fmt.Println("In NewDevCmd")
	cfg := config.ReadConfig()

	return &DevCmd{
		Config: cfg,
	}
}

func (d *DevCmd) Logs() {
	t, err := tail.TailFile(d.Config.LogFile, tail.Config{Follow: true})
	utils.ExitOnError(err)
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
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
	dh.StartContainerUI(cfg.DockerImages.Server, envVars, cfg.ApiServer.Port, containerName, []string{"ark", "server"}, mounts...)

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
