package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	ps "github.com/mitchellh/go-ps"
	"gopkg.in/yaml.v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Returnh true if current process is a child process of pulumi
func IsPulumiChild(args []string) bool {
	// Get parent pid
	pid := os.Getppid()
	proc, err := ps.FindProcess(pid)
	if err != nil {
		panic(err)
	}

	if proc == nil {
		return false
	}

	// Get binary name
	binName := proc.Executable()

	return strings.Contains(binName, "pulumi-")
}

func ReturnError(err error) error {
	if err != nil {
		return err
	}

	return nil
}

func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Marshals a struct of type 'V' to a yaml string
func YamlMarshall[T any](message T) (string, error) {
	// Convert message into yaml
	b, err := yaml.Marshal(message)
	if err != nil {
		fmt.Println("Could not covert request to yaml config data")
		log.Printf("yamlMarshall error: %v\n", message)
	}

	return string(b), err
}

// Unmarshals a string to the provided type 'V'
func JsonUnmarshall[T any](message string) (T, error) {
	var msg T
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		log.Println("Invalid message:" + err.Error())
		log.Println("jsonUnmarshall error:" + message)
	}
	return msg, err
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
