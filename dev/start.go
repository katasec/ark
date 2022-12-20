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

	fmt.Println("The resource group name:" + rgName)

	cfg := config.ReadConfig()
	cfg.AzureConfig.ResourceGroupName = rgName

	cfg.Save()

}
