/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/katasec/ark/cmd"
	"github.com/katasec/ark/devcmd"
	"github.com/katasec/ark/utils"
)

func main() {

	// Run Dev Init function if this program is being called by Pulumi
	if utils.IsPulumiChild(os.Args) {
		d := devcmd.NewDevCmd()
		d.Setup()
	}

	// Behave as normal cli
	cmd.Execute()
}
