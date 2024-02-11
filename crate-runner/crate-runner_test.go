package craterunner

import (
	"log"
	"testing"

	"github.com/katasec/ark/resources/azure/cloudspaces"
)

func TestTerraformApply(t *testing.T) {

	// Generate data for crate
	data := getAcsData()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-hello:v0.0.4", data)

	// Run crate
	r.Apply()
}

func TestTerraformDestroy(t *testing.T) {

	// Generate data for crate
	data := getAcsData()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-hello:v0.0.4", data)

	// Run crate
	r.Destroy()
}

func TestRunPulumiHelloApply(t *testing.T) {

	// Generate data for crate
	data := getAcsData()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-phello:v0.0.1", data)

	// Run crate
	r.Apply()
}

func TestRunPulumiHelloDestroy(t *testing.T) {

	// Generate data for crate
	acs := cloudspaces.NewAzureCloudSpace("UAE North")
	acs.AddSpoke("dev")
	data := acs.ToJson()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-phello:v0.0.1", data)

	// Run crate
	r.Destroy()
}

func TestPulumiAcsApply(t *testing.T) {

	// Generate data for crate
	data := getAcsData()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-azurecloudspace:v0.0.1", data)

	// Run crate
	r.Apply()
}

func TestPulumiAcsDestroy(t *testing.T) {

	// Generate data for crate
	data := getAcsData()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-azurecloudspace:v0.0.1", data)

	// Run crate
	r.Destroy()
}

func getAcsData() string {
	// Generate data for crate
	acs := cloudspaces.NewAzureCloudSpace("UAE North")
	acs.AddSpoke("dev")
	data := acs.ToJson()
	return data
}

func TestJson(t *testing.T) {
	data := getAcsData()
	log.Println(data)
}
