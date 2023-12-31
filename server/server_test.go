package server

import (
	"fmt"
	"testing"

	"github.com/katasec/ark/resources/azure/cloudspaces"

	_ "github.com/lib/pq" // Import the pq driver
)

func TestCreateAcsJson(t *testing.T) {

	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("prod")
	fmt.Println(acs.ToJson())

}

// CreateSampleAcs Creates a sample Azure Cloud Space Struct for testing
func CreateSampleAcs(name string) *cloudspaces.AzureCloudspace {

	return &cloudspaces.AzureCloudspace{
		Name: name,
		Hub: cloudspaces.VNETInfo{
			Name:          "hub",
			AddressPrefix: "10.1.0.0/16",
			SubnetsInfo: []cloudspaces.SubnetsInfo{
				{
					Name:          "hub-snet1",
					AddressPrefix: "10.1.1.0/24",
				},
			},
		},
		Spokes: []cloudspaces.VNETInfo{
			{
				Name:          "dev",
				AddressPrefix: "10.2.0.0/16",
				SubnetsInfo: []cloudspaces.SubnetsInfo{
					{
						Name:          "hub-snet1",
						AddressPrefix: "10.2.1.0/24",
					},
				},
			},
		},
	}
}
