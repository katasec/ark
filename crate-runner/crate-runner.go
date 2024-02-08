package craterunner

import (
	"log"
	"os"

	"github.com/katasec/ark/arkimage"
	"github.com/katasec/ark/prunner"
	"github.com/katasec/ark/tfrunner"
)

type CreateRunner struct {
	configData string
	imageName  string
}

func NewCrateRunner(imageName string, configData string) *CreateRunner {
	if !isValidImageName(imageName) {
		log.Println("Error validating image name")
		os.Exit(1)
	}

	return &CreateRunner{
		configData: configData,
		imageName:  imageName,
	}
}

func (c *CreateRunner) Run() {

	// Need to pull image to determine type
	image := arkimage.NewArkImage()
	imagetype := image.Pull(c.imageName)
	log.Println("The image type was:" + imagetype)

	switch imagetype {
	case "terraform":
		runner := tfrunner.NewTfrunner(c.imageName, c.configData)
		runner.Run()
	case "pulumi":
		runner := prunner.NewPRunner(c.imageName, c.configData)
		runner.Run()
	}
}
