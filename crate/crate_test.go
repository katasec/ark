package crate

import (
	"testing"
)

func TestPush(t *testing.T) {
	c := NewCrate()

	c.Push("https://github.com/katasec/ark-resource-azurecloudspace", "v0.0.1")
}

func TestPull(t *testing.T) {
	c := NewCrate()

	c.Pull("ghcr.io/katasec/azurecloudspace:v0.0.1")
}
