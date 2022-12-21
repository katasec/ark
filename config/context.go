package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/katasec/ark/utils"
	"gopkg.in/yaml.v2"
)

var (
	ArkDir  string
	cfgFile = fmt.Sprintf("%s/config", ArkDir)
)

type Config struct {
	Cloud         string `yaml:"cloud"`
	AzureConfig   AzureConfig
	StorageConfig StorageConfig
}

type AzureConfig struct {
	ResourceGroupName   string
	StorageConfig       StorageConfig
	ServiceBusNameSpace string
}

type StorageConfig struct {
	AccountName      string
	ConncetionString string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewConfig(cloudId string) *Config {

	// Define config dir
	ArkDir = getArkDir()

	// Create config directory
	if _, err := os.Stat(ArkDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(ArkDir, os.ModePerm)
		check(err)
	}

	// Create config yaml
	myConfig := &Config{
		Cloud: cloudId,
	}
	yamlData, err := yaml.Marshal(myConfig)
	check(err)

	// Save file
	cfgFile = fmt.Sprintf("%s/config", ArkDir)
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

	fmt.Println("Config was saved!")
}

func ReadConfig() *Config {

	configFile := getConfigFileName()

	// Create config directory
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		utils.ExitOnError(err)
	}

	bCfg, err := os.ReadFile(configFile)
	utils.ExitOnError(err)

	cfg := &Config{}
	yaml.Unmarshal(bCfg, cfg)

	return cfg
}

func getArkDir() string {
	// Define config dir
	homeDir, _ := os.UserHomeDir()
	ArkDir = homeDir + "/.ark"

	return ArkDir
}

func getConfigFileName() string {
	return fmt.Sprintf("%s/config", getArkDir())
}

func (cfg *Config) Dump() {
	yamlData, err := yaml.Marshal(cfg)
	utils.ExitOnError(err)

	fmt.Println(string(yamlData))
}
