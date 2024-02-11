package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/katasec/ark/resources/azure/cloudspaces"
	"github.com/katasec/ark/utils"
)

func TestAcsJson(t *testing.T) {
	configdata := genAcsJson()
	log.Println(configdata)
}
func genAcsJson() string {
	// Generate data for crate
	acs := cloudspaces.NewAzureCloudSpace("UAE North")
	acs.AddSpoke("dev")
	data := acs.ToJson()
	return data
}

func TestRandomString(t *testing.T) {
	x := utils.RandomString(6)
	fmt.Println(x)
}
