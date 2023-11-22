/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	dev "github.com/katasec/ark/cmd/dev"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup dependencies, for e.g. state store on Azure",
	Long:  "Setup dependencies, for e.g. state store on Azure",
	Run: func(cmd *cobra.Command, args []string) {
		d := dev.NewDevCmd()
		d.Setup()
	},
}

func init() {
	devCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
