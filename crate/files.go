package crate

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func listFilesRecursively(dirPath string) []string {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// Check error walking path
		if err != nil {
			log.Println("Error walking path:", err)
			os.Exit(1)
		}

		// Add file to files if it is not a directory
		if !strings.Contains(path, ".git/") && !strings.Contains(path, ".gitignore") && !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Println("Error walking path:", err)
		os.Exit(1)
	}

	return files
}

// createTempDir creates a temp dir
func createTempDir() string {

	// Create a folder called ark in the system temp dir if it does not exist.
	tmpdirBase := filepath.Join(os.TempDir(), "ark")
	err := os.Mkdir(tmpdirBase, os.FileMode(0777))
	log.Println("Ark's base temp dir: " + tmpdirBase)

	// Exit on error
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		fmt.Println("could not create tmpdirBase, exitting." + tmpdirBase)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create a temp dir under ark for use
	tmpdir, err := os.MkdirTemp(tmpdirBase, "ark-remote")
	if err != nil {
		fmt.Println("could not create tmpdir, exitting.")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Return the temp dir
	return tmpdir
}
