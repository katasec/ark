package worker

import (
	"fmt"
	"log"

	"github.com/katasec/ark/tfrunner"
)

// terraformHandler Creates a pulumi program and injects the message as pulumi config
func (w *Worker) terraformHandler(action string, resourceName string, configdata string, c chan error) {

	log.Println("In terraform handler")
	//fmt.Println("The config was:" + configdata)
	fmt.Println("The resource name was:" + resourceName)
	fmt.Println("The action name was:" + action)
	imageName := fmt.Sprintf("ark-resource-%s:v0.0.3", resourceName)
	runner := tfrunner.NewTfrunner(imageName, configdata)
	runner.Run()

}
