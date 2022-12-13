package utils

import (
	"os"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

// Returnh true if current process is a child process of pulumi
func IsPulumiChild(args []string) bool {

	// Get parent pid
	pid := os.Getppid()
	proc, err := ps.FindProcess(pid)
	if err != nil {
		panic(err)
	}

	// Get binary name
	binName := proc.Executable()

	return strings.Contains(binName, "pulumi-")
}
