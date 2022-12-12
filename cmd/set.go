/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/katasec/ark/config"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set the cloudid to 'azure' or 'aws'",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cloudId != "aws" && cloudId != "azure" {
			fmt.Println("Error: cloudid must be 'azure' or 'aws'")
			os.Exit(1)
		}

		config.NewConfig(cloudId)
	},
}

var cloudId string

func init() {
	configCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&cloudId, "cloudid", "c", "", "Use 'azure' or 'aws'")
	setCmd.MarkFlagRequired("cloudid")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
