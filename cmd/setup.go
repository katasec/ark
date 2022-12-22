/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/katasec/ark/devcmd"
	"github.com/spf13/cobra"
)

// setupCmd represents the init command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cloud dependencies required for spinning up a local ark environment",
	Long:  "Setup cloud dependencies required for spinning up a local ark environment.",
	Run: func(cmd *cobra.Command, args []string) {
		devcmd := devcmd.NewDevCmd()
		devcmd.Setup()
	},
}

func init() {
	devCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
