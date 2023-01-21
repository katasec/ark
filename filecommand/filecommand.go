package filecommand

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/katasec/ark/config"
	"gopkg.in/yaml.v2"
)

var (
	arkConfig = config.ReadConfig()
)

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
