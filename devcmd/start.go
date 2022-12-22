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

	orgName, err := getDefaultPulumiOrg()
	utils.ExitOnError(err)

	// Prefix org name with stack for FQDN
	stackFQDN := fmt.Sprintf("%s/ark-init/dev", orgName)

	// Get rg name
	rgName, err := getReference(stackFQDN, "rgName")
	utils.ExitOnError(err)

	// Get account name
	strgAcct, err := getReference(stackFQDN, "accountName")
	utils.ExitOnError(err)

	// Get namespace name
	ns, err := getReference(stackFQDN, "ns")
	utils.ExitOnError(err)

	// Write details to config file
	cfg := config.ReadConfig()
	cfg.AzureConfig.ResourceGroupName = rgName
	cfg.AzureConfig.ServiceBusNameSpace = ns
	cfg.AzureConfig.StorageConfig.AccountName = strgAcct
	cfg.Save()

	// Output config file
	//cfg.Dump()
}
