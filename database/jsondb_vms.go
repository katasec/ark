package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	resources "github.com/katasec/ark/resources"
)

// GetVm returns a vm if it exists
func (j *JsonRepository) GetVm(name string) (resources.Vm, error) {
	for _, vm := range j.vms {
		if strings.EqualFold(vm.Name, name) {
			return vm, nil
		}
	}
	return resources.Vm{}, errors.New("VM not found")
}

// AddVm Adds Vm
func (j *JsonRepository) AddVm(vm resources.Vm) resources.Vm {

	var myvm resources.Vm

	// Add if VM doesn't exist
	myvm, err := j.GetVm(vm.Name)
	if err != nil {
		j.vms = append(j.vms, vm)
		fmt.Println("Added VM")
	}

	return myvm
}

// SaveVms saves VMs to local file
func (j *JsonRepository) SaveVms() {

	fmt.Println("Saving VM")

	// Pretty print json
	jsonData, err := json.MarshalIndent(j.vms, "", "  ")
	fmt.Println(string(jsonData))
	logError(err)

	// Get file name
	f := j.OpenFile(JsonFile.Vms)
	fmt.Println("FileName:" + JsonFile.Vms)

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

func (j *JsonRepository) ListVms() []resources.Vm {
	return j.vms
}
