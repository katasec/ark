package tfrunner

import (
	"testing"

	"github.com/katasec/ark/resources/azure/cloudspaces"
)

func TestRunHello(t *testing.T) {

	acs := cloudspaces.NewAzureCloudSpace()
	acs.AddSpoke("dev")
	acs.AddSpoke("prod")
	runner := NewTfrunner("ark-resource-hello:v0.0.3", acs.ToJson())
	runner.Run()
}
