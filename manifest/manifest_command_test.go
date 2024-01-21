package manifest

import (
	"fmt"
	"testing"
)

func TestCreateManifestCommand(t *testing.T) {

	x := NewManifestCommand("apply", "../scripts/yaml/cloudspace.yaml")

	fmt.Println("The action is: " + x.action)
	fmt.Println("The endpoint is: " + x.EndPoint)

	//fmt.Println("The manifest is:\n" + x.Manifest)
}
