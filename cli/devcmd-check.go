package cli

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/katasec/ark/shell"
	"github.com/katasec/ark/utils"
)

// var (
// 	checksPassed = true
// )

func CheckSetupPreReqs() bool {

	// Setup spinner
	spinner := utils.NewArkSpinner()

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

	// if !checkDockerStarted() {
	// 	checksPassed = false
	// }

	fmt.Println()
	if !checksPassed {
		spinner.ErrorStatusEvent("One or more of the above checks failed. Please correct and try again.")
	} else {
		spinner.SuccessStatusEvent("Pre-flight checks passed!")
	}

	return checksPassed
}

func isInstalled(cmd string) bool {

	status := false

	// Setup spinner
	spinner := utils.NewArkSpinner()

	// Show status
	note := fmt.Sprintf("Verify %s is installed.", cmd)
	spinner.Start(note)
	_, err := exec.LookPath(cmd)
	if err == nil {
		status = true
	}
	spinner.Stop(err, note)

	return status
}

func checkDockerStarted() bool {

	// Setup spinner
	spinner := utils.NewArkSpinner()
	note := "Check docker is started"
	spinner.Start(note)

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
	if err == nil {
		status = true
	}

	spinner.Stop(err, note)

	return status
}

func isAzLoggedIn() bool {

	status := false

	// Setup spinner
	spinner := utils.NewArkSpinner()
	note := "Verify az is logged in"
	spinner.Start(note)

	shellCmd := "az ad signed-in-user show --query userPrincipalName -o tsv"
	_, err := shell.ExecShellCmd(shellCmd)
	if err == nil {
		status = true
	}

	spinner.Stop(err, note)

	return status
}
