package arkimage

import (
	"testing"
)

func TestPush(t *testing.T) {
	image := NewArkImage()

	image.Push("https://github.com/katasec/ark-resource-azurecloudspace", "v0.0.1")
}

func TestPull(t *testing.T) {
	image := NewArkImage()

	image.Pull("ghcr.io/katasec/ark-resource-hello:v0.0.1")
}