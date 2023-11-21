/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"strings"

	"github.com/katasec/ark/cmd"
	"github.com/mitchellh/go-ps"
)

func main() {

	// Behave as normal cli
	cmd.Execute()
}

// Returnh true if current process is a child process of pulumi
func IsPulumiChild(args []string) bool {
	// Get parent pid
	pid := os.Getppid()
	proc, err := ps.FindProcess(pid)
	if err != nil {
		panic(err)
	}

	if proc == nil {
		return false
	}

	// Get binary name
	binName := proc.Executable()

	return strings.Contains(binName, "pulumi-")
}
