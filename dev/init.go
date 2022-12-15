package dev

import (
	"context"
	"os"

	"github.com/dapr/cli/pkg/print"
	pulumihelper "github.com/katasec/pulumi-helper"
)

func Create() {

	print.PendingStatusEvent(os.Stdout, "Creating cloud resources for local dev environment.")

	args := &pulumihelper.PulumiRunRemoteParameters{
		ProjectName: "ArkInit",
		StackName:   "dev",
		Destroy:     false,
		Plugins: []map[string]string{
			{
				"name":    "azure-native",
				"version": "v1.89.1",
			},
		},
		GitURL:      "https://github.com/katasec/ArkInit.git",
		ProjectPath: "Azure",
		Config: []map[string]string{
			{
				"name":  "azure-native:location",
				"value": "westus2",
			},
		},
		Runtime: "dotnet",
		Branch:  "refs/remotes/origin/main",
	}

	ctx := context.Background()
	err := pulumihelper.RunPulumiRemote(ctx, args)

	if err != nil {
		print.FailureStatusEvent(os.Stdout, "Cloud not complete creation of cloud resources.")
	} else {
		print.SuccessStatusEvent(os.Stdout, "Created cloud resources for local dev environment")
	}

}

func Delete() {

	args := &pulumihelper.PulumiRunRemoteParameters{
		ProjectName: "ArkInit",
		StackName:   "dev",
		Destroy:     true,
		Plugins: []map[string]string{
			{
				"name":    "azure-native",
				"version": "v1.89.1",
			},
		},
		GitURL:      "https://github.com/katasec/ArkInit.git",
		ProjectPath: "Azure",
		Config: []map[string]string{
			{
				"name":  "azure-native:location",
				"value": "westus2",
			},
		},
		Runtime: "dotnet",
		Branch:  "refs/remotes/origin/main",
	}

	ctx := context.Background()
	pulumihelper.RunPulumiRemote(ctx, args)
}
