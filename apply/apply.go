package apply

import (
	"fmt"

	"github.com/katasec/ark/filecommand"
	"github.com/katasec/ark/utils"
)

func DoStuff(fileName string) {

	// Exit if file doesn't exist
	data := filecommand.ReadFile(fileName)

	// Get resource name
	resource, _ := filecommand.GetResource(data)

	// The kind argument in the file specified the resource
	// user wants to create
	kind := resource.Kind

	// Setup spinner
	s := utils.NewArkSpinner()
	msg := fmt.Sprintf("Starting apply for: %s\n", resource.Kind)
	s.InfoStatusEvent(msg)

	// fmt.Printf("Starting apply for: %s\n", resource.Kind)

	// Convert request to yaml for the API Server
	request, jsonContent, err := filecommand.Yaml2json(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch kind {
	case "azure/cloudspace":
		filecommand.CreateCloudspace(request, string(data))
		//filecommand.CreateCloudspace(request, jsonContent)
		fmt.Println(jsonContent)
	default:
		fmt.Println("Didn't recognize request")
	}

}
