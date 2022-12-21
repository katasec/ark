package dev

import (
	"fmt"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/utils"
)

func Start() {

	orgName, err := getDefaultPulumiOrg()
	utils.ExitOnError(err)

	stackFQDN := fmt.Sprintf("%s/ark-init/resource-group", orgName)
	rgName, err := getReference(stackFQDN, "rgName")
	utils.ExitOnError(err)

	stackFQDN = fmt.Sprintf("%s/ark-init/storage-account", orgName)
	strgAcct, err := getReference(stackFQDN, "accountName")
	utils.ExitOnError(err)

	stackFQDN = fmt.Sprintf("%s/ark-init/service-bus", orgName)
	ns, err := getReference(stackFQDN, "ns")
	utils.ExitOnError(err)

	cfg := config.ReadConfig()
	cfg.AzureConfig.ResourceGroupName = rgName
	cfg.AzureConfig.ServiceBusNameSpace = ns
	cfg.AzureConfig.StorageConfig.AccountName = strgAcct

	cfg.Save()

	cfg.Dump()
}
