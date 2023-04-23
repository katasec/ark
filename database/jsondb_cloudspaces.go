package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	resources "github.com/katasec/ark/resources/v0"
)

// GetCloudSpace returns a cloudspace if it exists
func (j *JsonRepository) GetCloudSpace(name string) (resources.CloudSpace, error) {
	for _, cs := range j.cloudspaces {
		if strings.EqualFold(cs.Name, name) {
			return cs, nil
		}
	}
	return resources.CloudSpace{}, errors.New("cloudspace not found")
}

// AddCloudSpace Adds Cloudspace
func (j *JsonRepository) AddCloudSpace(cs resources.CloudSpace) resources.CloudSpace {

	var item resources.CloudSpace

	// Add if VM doesn't exist
	fmt.Println("Looking for cloudspace:" + cs.Name)
	item, err := j.GetCloudSpace(cs.Name)

	if err != nil {
		j.cloudspaces = append(j.cloudspaces, cs)
		fmt.Println("Added Cloudpsace")
	}

	return item
}

// SaveCloudSpaces saves cloudspaces to local file
func (j *JsonRepository) SaveCloudSpaces() {

	fmt.Println("Saving Cloudpsace")

	// Pretty print json
	jsonData, err := json.MarshalIndent(j.cloudspaces, "", "  ")
	logError(err)

	// Get file name
	f := j.OpenFile(JsonFile.CloudSpaces)

	// Save json to file
	_, err = fmt.Fprintln(f, string(jsonData))
	if err == nil {
		log.Println("Info: Saved!")
	} else {
		log.Println("Error saving Vms")
		logError(err)
	}

	defer f.Close()
}
