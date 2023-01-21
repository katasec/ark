package apply

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/katasec/ark/filecommand"
)

func DoStuff(fileName string) {

	// Exit if file doesn't exist
	data := readFile(fileName)

	// Get resource name
	resource, _ := filecommand.GetResource(data)

	// The kind argument in the file specified the resource
	// user wants to create
	kind := resource.Kind
	fmt.Printf("Starting apply for: %s\n", resource.Kind)

	// Convert request to yaml for the API Server
	request, jsonContent, _ := filecommand.Yaml2json(data)

	switch kind {
	case "azure/cloudspace":
		createCloudspace(request, jsonContent)
	default:
		fmt.Println("Didn't recognize request")
	}

}

func createCloudspace(request filecommand.Cloudspace, jsonContent string) {

	// Get API endpoint
	endpoint := filecommand.GetApiEndpoint(request)

	// Send request to endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(jsonContent)))
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Output response and status
	fmt.Println(resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
