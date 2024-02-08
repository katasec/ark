package arkimage

import (
	"testing"
)

func TestPushHello(t *testing.T) {
	image := NewArkImage()
	image.Push("https://github.com/katasec/ark-resource-hello", "v0.0.4", "terraform")
}

func TestPushPHello(t *testing.T) {
	image := NewArkImage()
	image.Push("https://github.com/katasec/ark-resource-phello", "v0.0.1", "pulumi")
}

func TestPushAzureCloudSpace(t *testing.T) {
	image := NewArkImage()
	image.Push("https://github.com/katasec/ark-resource-azurecloudspace", "v0.0.1", "pulumi")
}

func TestPullHello(t *testing.T) {
	image := NewArkImage()
	image.Pull("ghcr.io/katasec/ark-resource-hello:v0.0.4")
}

func TestPullAzureCloudSpace(t *testing.T) {
	image := NewArkImage()
	image.Pull("ghcr.io/katasec/ark-resource-azurecloudspace:v0.0.1")
}
