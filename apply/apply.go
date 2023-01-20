package apply

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

	"github.com/katasec/ark/config"
	"gopkg.in/yaml.v2"
)

func DoStuff(fileName string) {

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
	request, jsonContent, _ := yaml2json(data)
	fmt.Println(request.Resource)

	// post data to apiserver
	arkConfig := config.ReadConfig()

	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%s/", arkConfig.ApiServer.Host, arkConfig.ApiServer.Port))
	endpoint.Path = request.Resource
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
	fmt.Println("Apply Started")

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
