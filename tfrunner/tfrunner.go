package tfrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/katasec/ark/arkimage"
)

type Tfrunner struct {
	ArkImage   string
	ExecPath   string
	configdata string
}

func NewTfrunner(arkImage string, configdata string) *Tfrunner {

	runner := &Tfrunner{
		ArkImage:   arkImage,
		configdata: configdata,
	}

	runner.installTerraform()

	return runner
}

func (t *Tfrunner) installTerraform() {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.6.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	t.ExecPath = execPath
}
func (t *Tfrunner) Run() {
	c := arkimage.NewArkImage()
	c.Pull(t.ArkImage)

	// Set working directory
	workingDir := c.GetLocalPath(t.ArkImage)
	log.Println("Working directory: " + workingDir)
	tf, err := tfexec.NewTerraform(workingDir, t.ExecPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	} else {
		log.Println("Terraform initialized")
	}

	// Unmarshal configdata
	var data map[string]interface{}
	err = json.Unmarshal([]byte(t.configdata), &data)
	if err != nil {
		panic(err)
	}

	// Marshal with indentation
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	// Create arkdata.json file
	arkdataFile, err := os.Create(workingDir + "/configdata.json")
	if err != nil {
		log.Fatalf("error creating configdata.json: %s", err)
	}
	defer arkdataFile.Close()
	arkdataFile.Write(prettyJSON)

	// Set stdout/stderr for TF output
	var buffer1 bytes.Buffer
	writer := io.MultiWriter(os.Stdout, &buffer1)
	tf.SetStdout(writer)
	tf.SetStderr(os.Stderr)

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	tf.Apply(context.Background())
	tf.Destroy(context.Background())
	fmt.Println(state.FormatVersion) // "0.1"
}
