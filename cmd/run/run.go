package run

import (
	"context"
	"fmt"
	"os"
	"strings"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func DoStuff(ics string) {
	//ics = "ghcr.io/katasec/cloudspace:v1"
	refx := strings.Split(ics, ":")[0]
	tagx := strings.Split(ics, ":")[1]
	regx := strings.Split(ics, "/")[0]

	fmt.Println("Doing stuff:" + ics)
	localpath := "/Users/ameerdeen/.ark/manifests/"
	fs, err := file.New(localpath)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// 1. Connect to a remote repository
	ctx := context.Background()
	//reg := "ghcr.io"
	//repo, err := remote.NewRepository(reg + "/katasec/artifact")
	repo, err := remote.NewRepository(refx)
	if err != nil {
		panic(err)
	}

	// Note: The below code can be omitted if authentication is not required
	username := os.Getenv("ARK_REGISTRY_USERNAME")
	password := os.Getenv("ARK_REGISTRY_PASSWORD")

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: auth.StaticCredential(regx, auth.Credential{
			Username: username,
			Password: password,
		}),
	}

	//tag := "v1"
	_, err = oras.Copy(ctx, repo, tagx, fs, tagx, oras.DefaultCopyOptions)
	//oras.FetchBytes(ctx, repo, tagx)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(repo.Reference.Repository + ":" + tagx + " copied to " + localpath)
	}

}
