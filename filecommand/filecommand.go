package filecommand

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/katasec/ark/config"
	"gopkg.in/yaml.v2"
)

var (
	arkConfig *config.Config
	dtLayout  string
)

func init() {
	//fmt.Println("In file command")
	arkConfig = config.ReadConfig()
	dtLayout = "2006-01-02 15:04:05"
}

func GetResource(data []byte) (Resource, error) {
	// convert to struct
	request := Resource{}
	err := yaml.Unmarshal(data, &request)
	if err != nil {
		log.Println(err)
	}

	return request, err
}

func GetApiEndpoint(request Cloudspace) string {
	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%s/", arkConfig.ApiServer.Host, arkConfig.ApiServer.Port))
	endpoint.Path = request.Kind
	if err != nil {
		log.Fatal(err)
	}

	return endpoint.String()
}

func Yaml2json(content []byte) (Cloudspace, string, error) {

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

func CreateCloudspace(request Cloudspace, fileContent string, method ...string) {

	var httpMethod string

	if len(method) == 0 {
		httpMethod = "post"
	} else {
		httpMethod = method[0]
	}
	httpMethod = strings.ToUpper(httpMethod)

	// Get API endpoint
	endpoint := GetApiEndpoint(request)

	// Send request to endpoint
	if httpMethod == "POST" {
		// Send post request
		resp, err := http.Post(endpoint, "application/x-yaml", bytes.NewBuffer([]byte(fileContent)))
		if err != nil {
			fmt.Println(err.Error())
		}

		// Output response and status
		if resp == nil {
			fmt.Println("Error, response is nil, exitting...")
			os.Exit(0)
		}
		fmt.Println(resp.Status)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(body))

		resp.Body.Close()

	} else {
		// create delete request
		req, err := http.NewRequest("DELETE", endpoint, bytes.NewBuffer([]byte(fileContent)))
		req.Header.Set("Content-Type", "application/x-yaml")
		if err != nil {
			fmt.Println(err.Error())
		}

		// Create client
		client := &http.Client{}

		// Fetch Request
		fmt.Println("Deleting:", fileContent)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("Delete:", resp.Status)
		}

		// // Display Results
		// fmt.Println("response Status : ", resp.Status)
		// fmt.Println("response Headers : ", resp.Header)
		// fmt.Println("response Body : ", string(respBody))
	}

}

func ReadFile(fileName string) []byte {
	// Exit if file doesn't exist
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		log.Println(fileName + " does not exist")
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
