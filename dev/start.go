package dev

import (
	"fmt"

	"github.com/katasec/ark/utils"
)

func Start() {
	fmt.Println("Ok - let's start")

	orgName, err := getDefaultPulumiOrg()
	utils.ExitOnError(err)

	stackFQDN := fmt.Sprintf("%s/ark-init/resource-group", orgName)
	rgName, err := getReference(stackFQDN, "rgName")
	utils.ExitOnError(err)

	fmt.Println("The resource group name:" + rgName)

}
