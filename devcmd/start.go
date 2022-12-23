package devcmd

import (
	"fmt"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/utils"
)

func Start() {
	refreshConfig()
}

func refreshConfig() {

	// Get config Details
	cfg := config.ReadConfig()

	// Contruct Pulumi Stack FQDN
	orgName, err := getDefaultPulumiOrg()
	utils.ExitOnError(err)
	stackFQDN := fmt.Sprintf("%s/ark-init/dev", orgName)

	/*
	 Extract Azure resource details from Pulumi Exports
	*/

	// Resource Group Name
	cfg.AzureConfig.ResourceGroupName, err = getReference(stackFQDN, ResourceGroupName)
	utils.ExitOnError(err)

	// Mq Connection String
	cfg.AzureConfig.MqConfig.MqConnectionString, err = getReference(stackFQDN, MqConnectionString)
	utils.ExitOnError(err)

	// Mq Connection String
	cfg.AzureConfig.MqConfig.MqName, err = getReference(stackFQDN, CommandQueueName)
	utils.ExitOnError(err)

	// Log Storage Account Endpoint
	cfg.AzureConfig.StorageConfig.LogStorageEndpoint, err = getReference(stackFQDN, LogStorageEndpoint)
	utils.ExitOnError(err)

	// Log Storage LogStorageKey
	cfg.AzureConfig.StorageConfig.LogStorageKey, err = getReference(stackFQDN, LogStorageKey)
	utils.ExitOnError(err)

	// Log Storage Container
	cfg.AzureConfig.StorageConfig.LogsContainer, err = getReference(stackFQDN, LogContainerName)
	utils.ExitOnError(err)

	// Pulumi State Container
	cfg.AzureConfig.StorageConfig.PulumiStateContainer, err = getReference(stackFQDN, PulumiStateContainerName)
	utils.ExitOnError(err)

	// Save Azure resource details to config file
	cfg.Save()

	// Output config file
	cfg.Dump()
}
