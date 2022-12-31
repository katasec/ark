package config

import (
	"errors"
	"fmt"
	"os"
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

	// if _, err := os.Stat(ArkDir); errors.Is(err, os.ErrNotExist) {
	// 	err := os.Mkdir(ArkDir, os.ModePerm)
	// 	check(err)

	// 	err = os.Mkdir(GetDbDir(), os.ModePerm)
	// 	check(err)
	// }

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
	createDir(GetDbDir())

	// if _, err := os.Stat(ArkDir); errors.Is(err, os.ErrNotExist) {
	// 	err := os.Mkdir(ArkDir, os.ModePerm)
	// 	check(err)

	// 	fmt.Println("I am here")
	// 	err = os.Mkdir(GetDbDir(), os.ModePerm)
	// 	check(err)

	// }

	// Create config yaml
	myConfig := &Config{
		CloudId: cloudId,
		LogFile: filepath.Join(ArkDir, "ark.log"),
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
	yaml.Unmarshal(bCfg, cfg)

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
