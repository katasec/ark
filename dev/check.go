package dev

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/dapr/cli/pkg/print"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/katasec/ark/utils"
)

var (
	checksPassed = true
)

func CheckStuff() {

	if !isInstalled("az") {
		checksPassed = false
	}

	if !isInstalled("docker") {
		checksPassed = false
	}

	if !isInstalled("pulumi") {
		checksPassed = false
	}

	if !isAzLoggedIn() {
		checksPassed = false
	}

	if !checkDockerStarted() {
		checksPassed = false
	}

	fmt.Println()
	if !checksPassed {
		print.InfoStatusEvent(os.Stdout, "One or more of the above checks failed. Please correct and try again.")
	} else {
		print.SuccessStatusEvent(os.Stdout, "Pre-flight checks passed!")
	}
}

func isInstalled(cmd string) bool {

	status := false

	// Show status
	stopSpinning := print.Spinner(os.Stdout, fmt.Sprintf("Verify %s is installed.", cmd))

	_, err := exec.LookPath(cmd)
	if err == nil {
		status = true
		stopSpinning(print.Success)
	} else {
		stopSpinning(print.Failure)
	}

	return status
}

func checkDockerStarted() bool {

	status := false
	ctx := context.Background()

	// Create Docker Client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Get list of containers
	_, err = cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})

	if err != nil {
		print.FailureStatusEvent(os.Stdout, "Please start docker for running Ark dev")
	} else {
		print.SuccessStatusEvent(os.Stdout, "Docker is started !")
		status = true
	}

	return status
}

func isAzLoggedIn() bool {

	status := false
	stopSpinning := print.Spinner(os.Stdout, "Verify az is logged in")

	shellCmd := "az ad signed-in-user show --query userPrincipalName -o tsv"
	_, err := utils.ExecShellCmd(shellCmd)
	if err == nil {
		stopSpinning(print.Success)
		status = true
	} else {
		stopSpinning(print.Failure)
		status = false
	}

	return status
}
