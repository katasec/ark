package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

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

func ExecShellCmd2(shellCmdStr string) error {

	shellCmd := strings.Split(shellCmdStr, " ")

	var cmd *exec.Cmd

	if len(shellCmd) > 1 {
		// shell cmd has arguments
		cmd = exec.Command(shellCmd[0], shellCmd[1:]...)
	} else if len(shellCmd) == 1 {
		// shell cmd has no args
		cmd = exec.Command(shellCmd[0], shellCmd[1:]...)
	} else if len(shellCmd) == 0 {
		err := errors.New("no shell command was received")
		return err
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	fmt.Println(out.String())

	if err != nil {
		return errors.New(errOut.String())
	}
	return nil
}
