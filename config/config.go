package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/katasec/ark/utils"
	"gopkg.in/yaml.v2"
)

var (
	ArkDir    string
	cfgFile   = fmt.Sprintf("%s/config", ArkDir)
	apiServer = ApiServer{
		Host: "0.0.0.0",
		Port: "5067",
	}
	DEV_ARK_SERVER_IMAGE_NAME = "ghcr.io/katasec/arkserver:v0.0.3"
	DEV_ARK_WORKER_IMAGE_NAME = "ghcr.io/katasec/arkworker:v0.0.4"
)

type Config struct {
	CloudId      string
	AzureConfig  AzureConfig
	AwsConfig    AwsConfig
	LogFile      string
	ApiServer    ApiServer
	DockerImages DockerImages
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

	// Create config yaml
	myConfig := &Config{
		LogFile:   filepath.Join(ArkDir, "ark.log"),
		ApiServer: apiServer,
		DockerImages: DockerImages{
			Server: DEV_ARK_SERVER_IMAGE_NAME,
			Worker: DEV_ARK_WORKER_IMAGE_NAME,
		},
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

	// Create config yaml
	myConfig := &Config{
		CloudId:   cloudId,
		LogFile:   filepath.Join(ArkDir, "ark.log"),
		ApiServer: apiServer,
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
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(dir, 0777)
	}
}

func (cfg *Config) SetupDirectories() {
	createDir(GetArkDir())
	createDir(GetDbDir())
}
