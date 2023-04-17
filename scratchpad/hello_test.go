package somethingtest_test

import (
	"testing"

	resources "github.com/katasec/ark/resources/v0/azure/cloudspaces"
)

func TestAcs(t *testing.T) {
	acs := resources.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("uat")
	acs.AddSpoke("prod")

	t.Log("Hub Name:" + acs.Hub.Name)

	//fmt.Println("Hub Name:" + acs.Hub.Name)
	t.Log("Hub AddressPrefix:" + acs.Hub.AddressPrefix)

	for _, j := range acs.Spokes {
		t.Log("Spoke Name:" + j.Name)
		t.Log("Spoke AddressPrefix:" + j.AddressPrefix)
	}

}
