package apply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/katasec/ark/config"
	"gopkg.in/yaml.v2"
)

var (
	arkConfig = config.ReadConfig()
)

func DoStuff(fileName string) {

	// Exit if file doesn't exist
	data := readFile(fileName)

	// Get resource name
	resource, _ := getResource(data)

	// The kind argument in the file specified the resource
	// user wants to create
	kind := resource.Kind
	fmt.Printf("Starting apply for: %s\n", resource.Kind)

	// Convert request to yaml for the API Server
	request, jsonContent, _ := yaml2json(data)

	switch kind {
	case "azure/cloudspace":
		createCloudspace(request, jsonContent)
	default:
		fmt.Println("Didn't recognize request")
	}

}

func getResource(data []byte) (Resource, error) {
	// convert to struct
	request := Resource{}
	err := yaml.Unmarshal(data, &request)
	if err != nil {
		log.Println(err)
	}

	return request, err
}

func createCloudspace(request Cloudspace, jsonContent string) {
	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%s/", arkConfig.ApiServer.Host, arkConfig.ApiServer.Port))
	endpoint.Path = request.Kind
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(endpoint.String(), "application/json", bytes.NewBuffer([]byte(jsonContent)))
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func yaml2json(content []byte) (Cloudspace, string, error) {

	// convert to struct
	request := Cloudspace{}
	err := yaml.Unmarshal(content, &request)
	if err != nil {
		log.Println(err)
		return request, "", err
	}

	// convert to json
	requestBytes, err := json.Marshal(request)
	if err != nil {
		log.Println("Error converting string to json")
		return request, "", err
	}

	return request, string(requestBytes), nil
}
