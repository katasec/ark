package somethingtest_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/katasec/ark/repositories"
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
	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("bob")
	acs.AddSpoke("joe")

	t.Log(acs.ToJson())

}

func TestSTuff(t *testing.T) {

	// Create the DB and CLoudpace Table
	db := OpenDb()
	repo := repositories.NewAzureCloudSpaceRepository(db)
	repo.CreateTable(db)

	// Create a cloudspace
	acs := cloudspaces.NewAzureCloudSpace()
	fmt.Println(acs.ToJson())

	// Add cloudspace to the database
	repo.CreateCloudSpace(acs)
}

// func TestDb() {
// 	db := OpenDb()

// 	defer db.Close()

// 	repo := repositories.NewAzureCloudSpaceRepository(db)

// 	repo.CreateTable(db)
// 	acs := genCloudSpace()

// 	repo.CreateCloudSpace(acs)

// 	acs.Hub.Name = "test3"
// 	repo.UpdateCloudSpace(acs)

// 	//repo.DeleteCloudSpace(acs)

// }

func TestWriteThenReadCloudSpace(t *testing.T) {

	// Create DB
	db := OpenDb()
	defer db.Close()

	// Create CloudSpace
	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")

	// Create Repository and write CloudSpace to DB
	repo := repositories.NewAzureCloudSpaceRepository(db)
	_, err := repo.CreateCloudSpace(acs)
	if err != nil {
		fmt.Println(err)
	}

	// Read cloudspace from DB
	myacs, err := repo.GetCloudSpaces()
	if err != nil {
		fmt.Println(err)
	}

	// Output cloudspace to console
	for _, acs := range myacs {
		for _, spokes := range acs.Spokes {
			fmt.Println(spokes.Name)
		}
	}
}

func OpenDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
	}

	return db
}
