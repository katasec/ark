package worker

import "fmt"

// terraformHandler Creates a pulumi program and injects the message as pulumi config
func (w *Worker) terraformHandler(action string, resourceName string, configdata string, c chan error) {

	//fmt.Println("The config was:" + configdata)
	fmt.Println("The resource name was:" + resourceName)
	fmt.Println("The action name was:" + action)
	//runner := NewTfrunner("ark-resource-hello:v0.0.3", acs.ToJson())
	//runner.Run()

}
