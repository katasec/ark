package tfrunner

import (
	"fmt"
	"testing"

	"github.com/katasec/ark/resources/azure/cloudspaces"
)

func TestRunHelloCrate(t *testing.T) {

	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("prod")

	fmt.Println(acs.ToJson())

	/*
		fmt.Println("Testing image download")
		runner := NewTfrunner("ghcr.io/katasec/ark-resource-hello:v0.0.1", acs.ToJson())
		runner.Run()
	*/
}
