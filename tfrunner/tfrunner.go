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
	tf         *tfexec.Terraform
}

func NewTfrunner(arkImage string, configdata string) *Tfrunner {

	runner := &Tfrunner{
		ArkImage:   arkImage,
		configdata: configdata,
	}

	// Install Terraform
	runner.installTerraform()

	// Run Terraform Init to prepare for apply/destroy
	runner.tf = runner.Init()

	return runner
}

func (runner *Tfrunner) installTerraform() {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.6.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	runner.ExecPath = execPath
}

func (runner *Tfrunner) Init() *tfexec.Terraform {
	c := arkimage.NewArkImage()
	c.Pull(runner.ArkImage)

	// Set working directory
	workingDir := c.GetLocalPath(runner.ArkImage)
	log.Println("Working directory: " + workingDir)
	tf, err := tfexec.NewTerraform(workingDir, runner.ExecPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	} else {
		log.Println("Terraform initialized")
	}

	// Unmarshal configdata
	var data map[string]interface{}
	err = json.Unmarshal([]byte(runner.configdata), &data)
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

	// Run Tf init
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	fmt.Println(state.FormatVersion) // "0.1"

	return tf
}

func (runner *Tfrunner) Apply() {
	runner.tf.Apply(context.Background())
}

func (runner *Tfrunner) Destroy() {
	runner.tf.Destroy(context.Background())
}
