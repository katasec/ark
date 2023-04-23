package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/katasec/ark/config"
	resources "github.com/katasec/ark/resources/v0"
)

type JsonFileNames struct {
	Vms         string
	CloudSpaces string
}

var JsonFile = &JsonFileNames{
	Vms:         path.Join(config.GetDbDir(), "vms.json"),
	CloudSpaces: path.Join(config.GetDbDir(), "cloudspaces.json"),
}

type JsonRepository struct {
	cloudspaces []resources.CloudSpace
	vms         []resources.Vm
	dbFiles     map[string]string
}

func NewJsonRepository() *JsonRepository {

	repo := &JsonRepository{
		cloudspaces: []resources.CloudSpace{},
		vms:         []resources.Vm{},
	}

	repo.InitDb()

	return repo
}

func (j *JsonRepository) InitDb() {

	// Define DB file names
	j.dbFiles = map[string]string{
		"vm":         JsonFile.Vms,
		"cloudspace": JsonFile.CloudSpaces,
	}

	// Make DB folder
	createDir(config.GetDbDir())

	// Create files
	for _, fileName := range j.dbFiles {
		f := j.OpenFile(fileName)
		f.Close()
	}

}

func (j *JsonRepository) OpenFile(filename string) *os.File {
	var f *os.File

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Creating db file...")
		// Create file if not exist
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Opening db file...")
		// Else open file
		f, err = os.OpenFile(filename, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
	}

	return f
}

func logError(err error) {
	if err != nil {
		log.Printf("%+v\n", err.Error())
	}
}

func createDir(dir string) {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(dir, os.ModePerm)
	}
}
