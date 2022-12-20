package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	ArkDir      string
	ContextFile = fmt.Sprintf("%s/config", ArkDir)
)

type Config struct {
	Cloud       string `yaml:"cloud"`
	AzureConfig AzureConfig
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
	ContextFile = fmt.Sprintf("%s/config", ArkDir)
	err = os.WriteFile(ContextFile, yamlData, 0644)
	check(err)

	return myConfig
}

func SaveContext() {

}

func ReadConfig() {

}

func getArkDir() string {
	// Define config dir
	homeDir, _ := os.UserHomeDir()
	ArkDir = homeDir + "/.ark"

	return ArkDir
}
