/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/katasec/ark/describe"
	"github.com/spf13/cobra"
)

var cloudspaceName string

// cloudspaceCmd represents the cloudspace command
var cloudspaceCmd = &cobra.Command{
	Use:   "cloudspace",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cloudspaceName == "" {
			cloudspaceName = "default"
		}

		describe.Start(cloudspaceName)
	},
}

func init() {
	describeCmd.AddCommand(cloudspaceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloudspaceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloudspaceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cloudspaceCmd.Flags().StringVarP(&cloudspaceName, "cloudspace", "c", cloudspaceName, "Name of cloudspace")

}
