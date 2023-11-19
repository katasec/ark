package scratchpad

import (
	"fmt"
	"testing"

	"github.com/katasec/ark/resources/v0/azure/cloudspaces"

	_ "github.com/mattn/go-sqlite3"
)

func TestAcs(t *testing.T) {
	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("uat")
	acs.AddSpoke("prod")
	acs.AddSpoke("prod2")
	acs.AddSpoke("uat") //shouldn't add this one

	t.Log("Hub Name:" + acs.Hub.Name)
	t.Log("Hub AddressPrefix:" + acs.Hub.AddressPrefix + "\n")

	//fmt.Println()
	t.Log("Spoke Count:" + fmt.Sprintf("%d", (len(acs.Spokes))))
	t.Log("Spoke Info")
	for _, i := range acs.Spokes {
		t.Log("\t" + i.Name + ":" + i.AddressPrefix)
	}

}
