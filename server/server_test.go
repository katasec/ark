package server

import (
	"fmt"
	"log"
	"testing"

	"github.com/katasec/ark/repositories"
	"github.com/katasec/ark/resources/v0/azure/cloudspaces"

	_ "github.com/lib/pq" // Import the pq driver
)

func TestNewAzureCloudSpaceRepository(t *testing.T) {

	server := NewServer()
	db, err := server.GetDbConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	acsRepo := repositories.NewAzureCloudSpaceRepository(db)

	//acsRepo.DeleteCloudSpace()
	name := "default4"
	acs := CreateSampleAcs(name)
	acsRepo.CreateCloudSpace(acs)

	acsFromDb, err := acsRepo.GetCloudSpace(name)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("ACS Name:%s\n", acsFromDb.Name)

	//acsRepo.DeleteCloudSpace(name)

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
