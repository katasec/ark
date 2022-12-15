/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	dev "github.com/katasec/ark/dev"
	"github.com/spf13/cobra"
)

// createCmd represents the init command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create cloud dependencies required for spinning up a local ark environment",
	Long:  "Create cloud dependencies required for spinning up a local ark environment.",
	Run: func(cmd *cobra.Command, args []string) {
		dev.Create()
	},
}

func init() {
	devCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
