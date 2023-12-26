package run

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func DoStuff(ics string) {

	// Extract registry domain, repo and tag from ics
	// For e.g. ics = ghcr.io/katasec/cloudspace:v1
	// registryDomain = ghcr.io
	// repo = katasec/cloudspace
	// tag = v1
	ref := strings.Split(ics, ":")[0]
	tagx := strings.Split(ics, ":")[1]
	registryDomain := strings.Split(ics, "/")[0]

	// Get home directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Create local path under home directory
	localpath := path.Join(homedir, "./ark/manifests") //"/Users/ameerdeen/.ark/manifests/"

	// Create a file store in the local path
	fs, err := file.New(localpath)
	fs.AllowPathTraversalOnWrite = true
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// Connect to a remote repository
	ctx := context.Background()
	repo, err := remote.NewRepository(ref)
	if err != nil {
		panic(err)
	}

	// Use the default registry credentials
	username := os.Getenv("ARK_REGISTRY_USERNAME")
	password := os.Getenv("ARK_REGISTRY_PASSWORD")
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: auth.StaticCredential(registryDomain, auth.Credential{
			Username: username,
			Password: password,
		}),
	}

	// Pull the image to the local file store
	_, err = oras.Copy(ctx, repo, tagx, fs, tagx, oras.DefaultCopyOptions)
	//descriptorStr := descriptor.

	fmt.Println(repo.Reference)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(repo.Reference.Repository + ":" + tagx + " copied to " + localpath)
	}

}
