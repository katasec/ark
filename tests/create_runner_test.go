package tests

import (
	"testing"

	craterunner "github.com/katasec/ark/crate-runner"
	"github.com/katasec/ark/resources/azure/cloudspaces"
)

func TestPulumi(t *testing.T) {
	pulumiApply()
	pulumiDestroy()
}

func TestPulumiAcs(t *testing.T) {
	pulumiApplyAcs()
	pulumiDestroyAcs()
}

func pulumiApply() {

	// Create crate runner
	r := craterunner.NewCrateRunner("ghcr.io/katasec/ark-resource-phello:v0.0.1", "")

	// Run crate
	r.Apply()
}

func pulumiApplyAcs() {

	configdata := genAcsJson()
	// Create crate runner
	r := craterunner.NewCrateRunner("ghcr.io/katasec/ark-resource-azurecloudspace:v0.0.1", configdata)

	// Run crate
	r.Apply()
}

func pulumiDestroyAcs() {
	configdata := genAcsJson()
	// Create crate runner
	r := craterunner.NewCrateRunner("ghcr.io/katasec/ark-resource-azurecloudspace:v0.0.1", configdata)

	// Run crate
	r.Destroy()
}

func pulumiDestroy() {
	// Create crate runner
	r := craterunner.NewCrateRunner("ghcr.io/katasec/ark-resource-phello:v0.0.1", "")

	// Run crate
	r.Destroy()
}

func TestTerraform(t *testing.T) {
	terraformApply()
	terraformDestroy()
}

func terraformApply() {
	// Generate data for crate
	data := getAcsData()

	// Create crate runner
	r := craterunner.NewCrateRunner("ghcr.io/katasec/ark-resource-hello:v0.0.4", data)

	// Run crate
	r.Destroy()
}

func terraformDestroy() {
	// Generate data for crate
	data := getAcsData()

	// Create crate runner
	r := craterunner.NewCrateRunner("ghcr.io/katasec/ark-resource-hello:v0.0.4", data)

	// Run crate
	r.Apply()
}

func getAcsData() string {
	// Generate data for crate
	acs := cloudspaces.NewAzureCloudSpace("UAE North")
	acs.AddSpoke("dev")
	data := acs.ToJson()
	return data
}
