package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/katasec/ark/utils"
	shell "github.com/katasec/utils/shell"
	"gopkg.in/yaml.v2"
)

var (
	ArkDir    string
	cfgFile   = fmt.Sprintf("%s/config", ArkDir)
	apiServer = ApiServer{
		Host: "0.0.0.0",
		Port: "5067",
	}
	DEV_ARK_SERVER_IMAGE_NAME = "ghcr.io/katasec/arkserver:v0.0.6"
	DEV_ARK_WORKER_IMAGE_NAME = "ghcr.io/katasec/arkworker:v0.0.6"
)

type Config struct {
	CloudId         string
	AzureConfig     AzureConfig
	AwsConfig       AwsConfig
	LogFile         string
	ApiServer       ApiServer
	DockerImages    DockerImages
	PulumiDefultOrg string
	RedisUrl        string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewEmptyConfig() {
	// Define config dir
	ArkDir = GetArkDir()

	// Create config directory
	createDir(ArkDir)
	createDir(GetDbDir())

	// Get Pulumi Default Org
	org, err := GetDefaultPulumiOrg()
	utils.ExitOnError(err)

	// Create config yaml
	funcport := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if funcport != "" {
		apiServer.Port = funcport
	}

	myConfig := &Config{
		LogFile:   filepath.Join(ArkDir, "ark.log"),
		ApiServer: apiServer,
		DockerImages: DockerImages{
			Server: DEV_ARK_SERVER_IMAGE_NAME,
			Worker: DEV_ARK_WORKER_IMAGE_NAME,
		},
		PulumiDefultOrg: org,
	}

	// Convert to yaml
	yamlData, err := yaml.Marshal(myConfig)
	check(err)

	// Save file
	cfgFile = filepath.Join(ArkDir, "config")
	err = os.WriteFile(cfgFile, yamlData, 0644)
	check(err)
}

func NewConfig(cloudId string) *Config {

	// Define config dir
	ArkDir = GetArkDir()

	// Create config directory
	createDir(ArkDir)
	createDir(GetDbDir())

	// Get Pulumi Default Org
	org, err := GetDefaultPulumiOrg()
	utils.ExitOnError(err)

	// Create config yaml
	myConfig := &Config{
		CloudId:         cloudId,
		LogFile:         filepath.Join(ArkDir, "ark.log"),
		ApiServer:       apiServer,
		PulumiDefultOrg: org,
	}

	yamlData, err := yaml.Marshal(myConfig)
	check(err)

	// Save file
	cfgFile = filepath.Join(ArkDir, "config")
	err = os.WriteFile(cfgFile, yamlData, 0644)
	check(err)

	return myConfig
}

func (cfg *Config) Save() {
	yamlData, err := yaml.Marshal(cfg)
	utils.ExitOnError(err)

	// Save file
	cfgFile = getConfigFileName()
	err = os.WriteFile(cfgFile, yamlData, 0644)
	utils.ExitOnError(err)

	// Setup spinner
	s := utils.NewArkSpinner()
	s.InfoStatusEvent("Updated config !")
}

func ReadConfig() *Config {

	configFile := getConfigFileName()

	// Check confif file exists
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		NewEmptyConfig()
	}

	bCfg, err := os.ReadFile(configFile)
	utils.ExitOnError(err)

	cfg := &Config{}

	// Get home dir
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	yaml.Unmarshal(bCfg, cfg)

	// Override log file
	cfg.LogFile = path.Join(dirname, ".ark", "ark.log")

	return cfg
}

func GetArkDir() string {
	// Define config dir
	homeDir, _ := os.UserHomeDir()
	ArkDir = filepath.Join(homeDir, ".ark")

	return ArkDir
}

func (cfg *Config) GetArkDir() string {
	// Define config dir
	homeDir, _ := os.UserHomeDir()
	ArkDir = filepath.Join(homeDir, ".ark")

	return ArkDir
}

func (cfg *Config) GetDbDir() string {
	// Define config dir
	homeDir, _ := os.UserHomeDir()
	dbDir := filepath.Join(homeDir, ".ark", "db/")

	return dbDir
}

func GetDbDir() string {
	// Define config dir
	homeDir, _ := os.UserHomeDir()
	dbDir := filepath.Join(homeDir, ".ark", "db/")

	return dbDir
}

func getConfigFileName() string {
	return fmt.Sprintf("%s/config", GetArkDir())
}

func (cfg *Config) Dump() {
	yamlData, err := yaml.Marshal(cfg)
	utils.ExitOnError(err)

	fmt.Println(string(yamlData))
}

func createDir(dir string) {
	fmt.Println("creating:" + dir)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		// Create dir & setup perms
		os.Mkdir(dir, 0777)
		setUnixPerms()
	}
}

func setUnixPerms() {
	// On non-wndows platform, the umask prevents creating
	// world readable folders,so shelling out.
	if runtime.GOOS != "windows" {
		shell.ExecShellCmd("chmod -R 777 " + GetArkDir())
	}
}

func (cfg *Config) SetupDirectories() {
	createDir(GetArkDir())
	createDir(GetDbDir())
}

func GetDefaultPulumiOrg() (string, error) {

	value, err := shell.ExecShellCmd("pulumi org get-default")
	utils.ExitOnError(err)

	// If no default org was set, then set current user
	// as default org
	if strings.Contains(value, "No Default") {

		// Get pulumi user
		whoami, err := shell.ExecShellCmd("pulumi whoami")
		utils.ExitOnError(err)
		whoami = strings.TrimSpace(whoami)

		// Set as default org
		cmd := fmt.Sprintf("pulumi org set-default %s", whoami)
		shell.ExecShellCmd(cmd)
		value = whoami
	}

	value = strings.TrimSpace(value)

	return value, err
}
