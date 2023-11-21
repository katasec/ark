/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/katasec/ark/cmd/delete"
	"github.com/spf13/cobra"
)

var deleteFile string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Use to delete resource defined in the yaml file",
	Long:  "Use to delete resource defined in the yaml file",
	Run: func(cmd *cobra.Command, args []string) {
		delete.DoStuff(deleteFile)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	deleteCmd.Flags().StringVarP(&deleteFile, "filename", "f", deleteFile, "Name of yaml file with configuration to apply")
	deleteCmd.MarkFlagRequired("filename")
}
