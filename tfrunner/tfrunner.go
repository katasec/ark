package tfrunner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/katasec/ark/crate"
)

type Tfrunner struct {
	CrateImage string
	ExecPath   string
}

func NewTfrunner(crateImage string, arkdata string) *Tfrunner {

	runner := &Tfrunner{
		CrateImage: crateImage,
	}

	runner.installTerraform()

	return runner
}

func (t *Tfrunner) installTerraform() {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	t.ExecPath = execPath
}
func (t *Tfrunner) Run() {
	c := crate.NewCrate()
	c.Pull(t.CrateImage)

	// Get home directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Set working directory
	workingDir := path.Join(homedir, "./ark/manifests")
	tf, err := tfexec.NewTerraform(workingDir, t.ExecPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	} else {
		log.Println("Terraform initialized")
	}

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

	fmt.Println(state.FormatVersion) // "0.1"
}
