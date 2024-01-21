package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type ManifestCommand struct {
	action    string
	manifest  string
	kind      string
	EndPoint  string
	Method    string
	manifestb []byte
}

func NewManifestCommand(action string, fileName string) *ManifestCommand {

	// Read manifest in file
	manifestb := readFile(fileName)

	// Get resource kind from manifest
	resource := getResource(manifestb)

	// Return manifest command
	return &ManifestCommand{
		action:    action,
		manifestb: manifestb,
		manifest:  string(manifestb),
		kind:      resource.Kind,
		EndPoint:  getApiEndpoint(resource.Kind),
		Method:    getMethod(action),
	}
}

func readFile(fileName string) []byte {
	// Exit if file doesn't exist
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		log.Println("Error reading file:" + fileName + " does not exist")
		os.Exit(1)
	}

	// ready file
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(data))
	return data
}

func getApiEndpoint(kind string) string {
	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%s/", arkConfig.ApiServer.Host, arkConfig.ApiServer.Port))
	endpoint.Path = kind
	if err != nil {
		log.Fatal(err)
	}
	return endpoint.String()
}

func getMethod(action string) (method string) {

	switch action {
	case "apply":
		return "POST"
	case "destroy":
		return "DELETE"
	default:
		log.Println("Invalid action: " + action)
		os.Exit(1)
	}

	return
}

func (c *ManifestCommand) Execute() {

	// Create httpClient
	httpClient := &http.Client{}

	// create request
	req, err := http.NewRequest(c.Method, c.EndPoint, bytes.NewBuffer(c.manifestb))
	req.Header.Set("Content-Type", "application/x-yaml")
	if err != nil {
		fmt.Println(err.Error())
	}

	// log request
	// log.Println("Action:", c.action)
	// log.Println("Endpoint:", c.EndPoint)
	// log.Println("Manifest:", c.manifest)
	// log.Println("Method:", c.Method)

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
