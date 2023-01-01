package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	resources "github.com/katasec/ark/resources"
)

// AddCloudSpace Adds Cloudspace
func (j *JsonRepository) AddCloudSpace(cs resources.CloudSpace) resources.CloudSpace {

	var item resources.CloudSpace

	// Add if VM doesn't exist
	item, err := j.GetCloudSpace(cs.Name)
	if err != nil {
		j.cloudspaces = append(j.cloudspaces, cs)
		fmt.Println("Added Cloudpsace")
	}

	return item
}

// GetCloudSpace returns a cloudspace if it exists
func (j *JsonRepository) GetCloudSpace(name string) (resources.CloudSpace, error) {
	for _, cs := range j.cloudspaces {
		if strings.ToLower(cs.Name) == name {
			return cs, nil
		}
	}
	return resources.CloudSpace{}, errors.New("cloudspace not found")
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
	_, err = f.WriteString(string(jsonData))
	logError(err)
	if err == nil {
		log.Println("Info: Saved!")
	}
	defer f.Close()
}
