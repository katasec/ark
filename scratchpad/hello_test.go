package somethingtest_test

import (
	"fmt"
	"testing"

	resources "github.com/katasec/ark/resources/v0/azure/cloudspaces"
)

func TestAcs(t *testing.T) {
	acs := resources.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("uat")
	acs.AddSpoke("prod")

	fmt.Println("Hub Name:" + acs.Hub.Name)
	fmt.Println("Hub AddressPrefix:" + acs.Hub.AddressPrefix)

	for _, j := range acs.Spokes {
		fmt.Println("Spoke Name:" + j.Name)
		fmt.Println("Spoke AddressPrefix:" + j.AddressPrefix)
	}

}
