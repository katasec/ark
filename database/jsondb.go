package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/katasec/ark/config"
	resources "github.com/katasec/ark/resources"
)

type JsonRepository struct {
	cloudspaces []resources.CloudSpace
	vms         []resources.Vm
	vmFile      string
}

// func JsonRepository(c resources.CloudspaceRequest) {

// }

func NewJsonRepository() *JsonRepository {

	// dir := config.GetArkDir()
	// dbDir := path.Join(config.ArkDir, "db")
	vmFile := path.Join(config.GetDbDir(), "vm.json")

	return &JsonRepository{
		cloudspaces: []resources.CloudSpace{},
		vms:         []resources.Vm{},
		vmFile:      vmFile,
	}
}

func (j *JsonRepository) AddVm(vm resources.Vm) resources.Vm {
	j.vms = append(j.vms, vm)
	return vm
}

func (j *JsonRepository) SaveVms() error {
	jsonInBytes, err := json.MarshalIndent(j.vms, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.Create(j.vmFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(string(jsonInBytes))

	return nil
}

func (j *JsonRepository) GetVm(name string) (resources.Vm, error) {
	for _, vm := range j.vms {
		if strings.ToLower(vm.Name) == name {
			return vm, nil
		}
	}
	return resources.Vm{}, errors.New("VM not found")
}

func (j *JsonRepository) ListVms() []resources.Vm {
	return j.vms
}
