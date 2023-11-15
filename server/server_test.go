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
	name := "default3"
	acsRepo.CreateCloudSpace(GetSampleAcs(name))

	acs, err := acsRepo.GetCloudSpace(name)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("ACS Name:%s\n", acs.Name)
}

func GetSampleAcs(name string) *cloudspaces.AzureCloudspace {

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
