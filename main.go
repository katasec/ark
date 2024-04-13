/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/katasec/ark/cmd"
	"os"
)

func main() {

	for i := range 10 {
		fmt.Println(i)
	}

	os.Exit(0)
	// Behave as normal cli
	cmd.Execute()

}
