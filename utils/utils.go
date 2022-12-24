package utils

import (
	"fmt"
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

func ReturnError(err error) error {
	if err != nil {
		return err
	}

	return nil
}

func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
