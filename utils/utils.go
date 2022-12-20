package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
func ExecShellCmd(shellCmdStr string) (output string, err error) {

	shellCmd := strings.Split(shellCmdStr, " ")

	var cmd *exec.Cmd

	if len(shellCmd) > 1 {
		// shell cmd has arguments
		cmd = exec.Command(shellCmd[0], shellCmd[1:]...)
	} else if len(shellCmd) == 1 {
		// shell cmd has no args
		cmd = exec.Command(shellCmd[0], shellCmd[1:]...)
	} else if len(shellCmd) == 0 {
		err = errors.New("no shell command was received")
		output = ""
		return output, err
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err = cmd.Run()

	if err != nil {
		return "", errors.New(errOut.String())
	} else {
		return out.String(), nil
	}

}
