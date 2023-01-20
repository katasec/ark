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
	ArkDir  string
	cfgFile = fmt.Sprintf("%s/config", ArkDir)
)

type Config struct {
	CloudId     string
	AzureConfig AzureConfig
	AwsConfig   AwsConfig
	LogFile     string
	ApiServer   ApiServer
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
		LogFile: filepath.Join(ArkDir, "ark.log"),
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
		CloudId: cloudId,
		LogFile: filepath.Join(ArkDir, "ark.log"),
		ApiServer: ApiServer{
			Host: "localhost",
			Port: "5067",
		},
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
	fmt.Println("The Logfile was:" + cfg.LogFile)

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
		os.Mkdir(dir, os.ModePerm)
	}
}
