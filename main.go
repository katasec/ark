/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/katasec/ark/cmd"

// "github.com/katasec/ark/cmd"
// "github.com/katasec/ark/dev"
// "github.com/katasec/ark/utils"
func main() {

	// Run if this program is being called from Pulumi
	// if utils.IsPulumiChild(os.Args) {
	// 	d := dev.NewDevCmd()
	// 	d.Setup()
	// }

	// Behave as normal cli
	//fmt.Println("hi")
	cmd.Execute()
}
