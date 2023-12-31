package tfrunner

import (
	"fmt"
	"testing"
)

func TestRunHelloCrate(t *testing.T) {
	fmt.Println("Testing crate download")
	runner := NewTfrunner("ghcr.io/katasec/ark-resource-hello:v0.0.1", "arkdata")
	runner.Run()
}
