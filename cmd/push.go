/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/katasec/ark/cmd/push"
	"github.com/spf13/cobra"
)

var gitUrl string
var gitTag string

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if gitUrl == "" || gitTag == "" {
			cmd.Help()
			os.Exit(0)
		}
		//push.DoPush("https://github.com/katasec/ark-resource-azurecloudspace.git", "v0.0.1")
		push.DoPush(gitUrl, gitTag)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pushCmd.Flags().StringVarP(&gitUrl, "giturl", "g", "", "Git url for e.g https://github.com/katasec/ark-resource-azurecloudspace.git")
	pushCmd.Flags().StringVarP(&gitTag, "gitTag", "t", "", "Git tag for e.g v0.0.1")
}
