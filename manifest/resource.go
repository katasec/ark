package manifest

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Resource struct {
	Kind string
}

func getResource(data []byte) Resource {
	// convert to struct
	resource := Resource{}
	err := yaml.Unmarshal(data, &resource)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return resource
}
