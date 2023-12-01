package cloudspaces

import (
	"fmt"
	"os"
	"testing"

	"github.com/katasec/tableio"
)

func TestGenJson(t *testing.T) {
	// Generate a new cloudspace
	acs := NewAzureCloudSpace()
	// Add a spoke to the cloudspace
	acs.AddSpoke("test")
	// Generate the JSON representation of the cloudspace
	json := acs.ToJson()
	fmt.Println(json)
}

func TestAcs(t *testing.T) {
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	acsTable, err := tableio.NewTableIO[AzureCloudspace]("postgres", conn)
	if err != nil {
		fmt.Println(err)
	}

	acsTable.CreateTableIfNotExists(false)

	// acs := cloudspaces.NewAzureCloudSpace()
	// acsTable.Insert(*acs)

	result, _ := acsTable.All()

	for _, newAcs := range result {
		fmt.Println(newAcs.Name)
		fmt.Println(newAcs.Hub.SubnetsInfo[0].AddressPrefix)
	}
}
