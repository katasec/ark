package craterunner

import (
	"log"
	"os"

	"github.com/katasec/ark/arkimage"
	"github.com/katasec/ark/prunner"
	"github.com/katasec/ark/tfrunner"
)

type CreateRunner struct {
	configData string // Config data to inject into the IaC program before run
	imageName  string // container image name
	imageType  string // For e.g. Pulumi or Terraform
	workDir    string // Pulumi program will be downloaded into this local folder
}

func NewCrateRunner(imageName string, configData string) *CreateRunner {
	if !isValidImageName(imageName) {
		log.Println("Error validating image name")
		os.Exit(1)
	}

	// Need to pull image to determine type
	image := arkimage.NewArkImage()
	imageType, workDir := image.Pull(imageName)
	log.Println("The image type was:" + imageType)

	return &CreateRunner{
		configData: configData,
		imageName:  imageName,
		imageType:  imageType,
		workDir:    workDir,
	}
}

func (c *CreateRunner) Apply() {

	switch c.imageType {
	case "terraform":
		runner := tfrunner.NewTfrunner(c.imageName, c.configData)
		runner.Apply()
	case "pulumi":
		runner := prunner.NewPRunner(c.imageName, c.configData, c.workDir)
		runner.Up()
	}
}

func (c *CreateRunner) Destroy() {

	switch c.imageType {
	case "terraform":
		runner := tfrunner.NewTfrunner(c.imageName, c.configData)
		runner.Destroy()
	case "pulumi":
		runner := prunner.NewPRunner(c.imageName, c.configData, c.workDir)
		runner.Destroy()
	}
}
