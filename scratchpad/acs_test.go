package somethingtest_test

import (
	"strconv"
	"testing"

	resources "github.com/katasec/ark/resources/v0/azure/cloudspaces"
)

func TestAcs(t *testing.T) {
	acs := resources.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("uat")
	acs.AddSpoke("prod")
	acs.AddSpoke("prod2")
	acs.AddSpoke("uat") //shouldn't add this one

	t.Log("Spoke Count:" + strconv.Itoa(len(acs.Spokes)))
	t.Log("Hub Name:" + acs.Hub.Name)
	t.Log("Hub AddressPrefix:" + acs.Hub.AddressPrefix + "\n")
	for _, i := range acs.Hub.SubnetsInfo {
		t.Log("\tSubnet Name:" + i.Name)
		t.Log("\tSubnet AddressPrefix:" + i.AddressPrefix)
	}

	for _, j := range acs.Spokes {
		t.Log("\n")
		t.Log("Spoke Name:" + j.Name)
		t.Log("Spoke AddressPrefix:" + j.AddressPrefix)
		for _, k := range j.SubnetsInfo {
			t.Log("\t\tSubnet Name:" + k.Name)
			t.Log("\t\tSubnet AddressPrefix:" + k.AddressPrefix)
		}
	}

}

func TestAcsJson(t *testing.T) {
	acs := resources.NewAzureCloudSpace()
	acs.AddSpoke("bob")
	acs.AddSpoke("joe")

	t.Log(acs.ToJson())

}
