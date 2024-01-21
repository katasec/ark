/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/katasec/ark/manifest"
	"github.com/spf13/cobra"
)

var applyFile string

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Use to create resource defined in the yaml file",
	Long:  "Use to create resource defined in the yaml file",
	Run: func(cmd *cobra.Command, args []string) {
		//apply.StartApply(applyFile)

		myCmd := manifest.NewManifestCommand("apply", applyFile)
		myCmd.Execute()
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	applyCmd.Flags().StringVarP(&applyFile, "filename", "f", applyFile, "Name of yaml file with configuration to apply")
	applyCmd.MarkFlagRequired("filename")
}
