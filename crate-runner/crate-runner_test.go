package craterunner

import (
	"testing"

	"github.com/katasec/ark/resources/azure/cloudspaces"
)

func TestRunTerraformCrate(t *testing.T) {

	// Generate data for crate
	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("prod")
	data := acs.ToJson()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-hello:v0.0.4", data)

	// Run crate
	r.Run()
}

func TestRunPulumiHelloCrate(t *testing.T) {

	// Generate data for crate
	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("prod")
	data := acs.ToJson()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-phello:v0.0.1", data)

	// Run crate
	r.Run()
}

func TestRunPulumiAzureCloudspaceCrate(t *testing.T) {

	// Generate data for crate
	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("prod")
	data := acs.ToJson()

	// Create crate runner
	r := NewCrateRunner("ghcr.io/katasec/ark-resource-azurecloudspace:v0.0.1", data)

	// Run crate
	r.Run()
}
