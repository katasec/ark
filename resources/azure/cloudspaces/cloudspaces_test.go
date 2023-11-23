package cloudspaces

import (
	"fmt"
	"testing"
)

func TestGenJson(t *testing.T) {
	// Generate a new cloudspace
	acs := NewAzureCloudSpace()
	// Add a spoke to the cloudspace
	acs.AddSpoke("test")
	// Generate the JSON representation of the cloudspace
	json := acs.ToJson()
	fmt.Println(json)
}
