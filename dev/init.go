package dev

import (
	"context"

	pulumihelper "github.com/katasec/pulumi-helper"
)

func Create() {

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
	pulumihelper.RunPulumiRemote(ctx, args)
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
