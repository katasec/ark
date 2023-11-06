package delete

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
	msg := fmt.Sprintf("Starting delete for: %s\n", resource.Kind)
	s.InfoStatusEvent(msg)

	// fmt.Printf("Starting delete for: %s\n", resource.Kind)

	// Convert request to yaml for the API Server
	request, _, _ := filecommand.Yaml2json(data)

	switch kind {
	case "azure/cloudspace":
		filecommand.CreateCloudspace(request, string(data), "delete")
	default:
		fmt.Println("Didn't recognize request")
	}
}
