package tests

import (
	"testing"

	"github.com/katasec/ark/arkimage"
)

func TestPushHello(t *testing.T) {
	image := arkimage.NewArkImage()
	image.Push("https://github.com/katasec/ark-resource-hello", "v0.0.4", "terraform")
}

func TestPushPHello(t *testing.T) {
	image := arkimage.NewArkImage()
	image.Push("https://github.com/katasec/ark-resource-phello", "v0.0.1", "pulumi")
}

func TestPushAzureCloudSpace(t *testing.T) {
	image := arkimage.NewArkImage()
	image.Push("https://github.com/katasec/ark-resource-azurecloudspace", "v0.0.3", "pulumi")
}

func TestPullHello(t *testing.T) {
	image := arkimage.NewArkImage()
	image.Pull("ghcr.io/katasec/ark-resource-hello:v0.0.4")
}

func TestPullAzureCloudSpace(t *testing.T) {
	image := arkimage.NewArkImage()
	image.Pull("ghcr.io/katasec/ark-resource-azurecloudspace:v0.0.1")
}
