package apply

import (
	"errors"
	"log"
	"os"
)

func readFile(fileName string) []byte {
	// Exit if file doesn't exist
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		log.Println(fileName + " does not exist")
		os.Exit(1)
	}

	// ready file
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
