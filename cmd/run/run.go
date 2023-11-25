package run

import (
	"context"
	"fmt"
	"os"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func DoStuff(recipe string) {
	fmt.Println("Doing stuff:" + recipe)
	localpath := "/Users/ameerdeen/.ark/manifests/"
	fs, err := file.New(localpath)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// 1. Connect to a remote repository
	ctx := context.Background()
	reg := "ghcr.io"
	repo, err := remote.NewRepository(reg + "/katasec/artifact")
	if err != nil {
		panic(err)
	}

	// Note: The below code can be omitted if authentication is not required
	username := os.Getenv("ARK_REGISTRY_USERNAME")
	password := os.Getenv("ARK_REGISTRY_PASSWORD")

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: auth.StaticCredential(reg, auth.Credential{
			Username: username,
			Password: password,
		}),
	}

	tag := "v1"
	manifestDescriptor, err := oras.Copy(ctx, repo, tag, fs, tag, oras.DefaultCopyOptions)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(repo.Reference.Repository + ":" + tag + " copied to " + localpath)
	}
	fmt.Println("manifest descriptor:", manifestDescriptor)
}
