package craterunner

import (
	"log"
	"strings"
)

func isValidImageName(imageName string) bool {
	isvalid := true

	if !strings.Contains(imageName, ":") {
		log.Printf("Error, image name must contain image version, received: %s \n", imageName)
		isvalid = false
	}

	return isvalid
}
