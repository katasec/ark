package dev

import (
	"context"
	"fmt"

	pulumihelper "github.com/katasec/pulumi-helper"
)

func Init() {
	fmt.Println("Init called")

	// Setup Pulumi run parameters
	args := &pulumihelper.PulumiRunRemoteParameters{
		ProjectName: "Azure",
		StackName:   "dev",
		Destroy:     false,
		Plugins: []map[string]string{
			{
				"name":    "azure-native",
				"version": "v1.64.1",
			},
		},
		GitURL:      "https://github.com/katasec/ArkInit.git",
		ProjectPath: "Azure",
		Config: []map[string]string{
			{
				"location": "westus2",
			},
		},
	}

	ctx := context.Background()
	pulumihelper.RunPulumiRemote(ctx, args)
}
