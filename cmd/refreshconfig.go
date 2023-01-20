/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/katasec/ark/dev"
	"github.com/spf13/cobra"
)

// refreshconfigCmd represents the refreshconfig command
var refreshconfigCmd = &cobra.Command{
	Use:   "refreshconfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		d := dev.NewDevCmd()
		d.RefreshConfig()
	},
}

func init() {
	devCmd.AddCommand(refreshconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refreshconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// refreshconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
