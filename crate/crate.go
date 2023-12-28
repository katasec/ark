package crate

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/katasec/ark/config"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

var (
	arkConfig = config.ReadConfig()
)

type Crate struct {
	ctx              context.Context
	registryDomain   string
	registryUsername string
	registryPassword string
}

func NewCrate() *Crate {
	return &Crate{
		ctx:              context.Background(),
		registryDomain:   arkConfig.RegistryDomain,
		registryUsername: arkConfig.RegistryUsername,
		registryPassword: arkConfig.RegistryPassword,
	}
}

// Push Pushes code from a git repo to a registry
func (c *Crate) Push(gitUrl string, tag string) {
	fmt.Println("Pushing crate")

	if isValidGitUrl(gitUrl) {
		// Clone repo into a temp dir
		repoDir := cloneRemote(gitUrl, tag)
		os.Chdir(repoDir)

		// Push the cloned directory to the registry
		c.pushToRegistry(repoDir, tag, gitUrl)
	} else {
		fmt.Println("Not a git URL: " + gitUrl)
	}

}

// Pull Pulls code from a registry to a local directory
func (c *Crate) Pull(image string) {
	//ghcr.io/katasec/azurecloudspace:v0.0.1

	fmt.Println(image)
	fmt.Println("Pulling crate")

	// Get home directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Create local path under home directory
	localpath := path.Join(homedir, "./ark/manifests")

	// Create a file store in the local path
	fs, err := file.New(localpath)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// Connect to a remote repository
	ref := strings.Split(image, ":")[0]
	repo, err := remote.NewRepository(ref)
	if err != nil {
		panic(err)
	}
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: auth.StaticCredential(c.registryDomain, auth.Credential{
			Username: c.registryUsername,
			Password: c.registryPassword,
		}),
	}

	// Pull the image to the local file store
	tagx := strings.Split(image, ":")[1]
	_, err = oras.Copy(c.ctx, repo, tagx, fs, tagx, oras.DefaultCopyOptions)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(repo.Reference.Repository + ":" + tagx + " copied to " + localpath)
	}

}

// pushToRegistry Pushes files from a directory to a registry
func (c *Crate) pushToRegistry(dir string, tag string, gitUrl string) {

	// Pusing a file to a registry requires the creation of an ORA
	// local file store and a remote repository. The file store is a local directory
	// where the file is stored. The remote repository is the registry where the file is pushed to.

	// 0. Create a file store
	fs, err := file.New(dir)
	fs.AllowPathTraversalOnWrite = true
	log.Println("Creating file store: " + dir)
	if err != nil {
		fmt.Println("Error creating file store:", err)
		os.Exit(1)
	}
	defer fs.Close()
	ctx := context.Background()

	// 1. Add files to the file store
	mediaType := "application/vnd.test.file"
	fileNames := listFilesRecursively(dir)
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))

	for _, fullPath := range fileNames {
		// Trim path relative to filestore
		relativePath := strings.TrimPrefix(fullPath, dir+"/")
		log.Println("Adding file: " + relativePath)

		// Add the file to the file store
		fileDescriptor, err := fs.Add(ctx, relativePath, mediaType, "")
		if err != nil {
			fmt.Println("Error adding file to file store:", err)
			os.Exit(1)
		}

		// Add the file descriptor to the list of file descriptors
		fileDescriptors = append(fileDescriptors, fileDescriptor)
	}

	// 2. Pack the files and tag the packed manifest
	artifactType := "application/vnd.test.artifact"
	opts := oras.PackManifestOptions{
		Layers: fileDescriptors,
	}
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1_RC4, artifactType, opts)
	if err != nil {
		fmt.Println("Error packing manifest:", err)
		os.Exit(1)
	}

	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		fmt.Println("Error tagging manifest:", err)
		os.Exit(1)
	}

	// 3. Get Remote registry details
	resourceName, _ := hasValidResourceName(gitUrl)
	repoName := arkConfig.Registry + "/" + resourceName
	repoName = strings.Replace(repoName, "//", "/", -1)
	registryDomain := arkConfig.RegistryDomain

	// 4. Connect to the remote repository
	log.Println("Connecting to: " + repoName)
	repo, err := remote.NewRepository(repoName)
	if err != nil {
		panic(err)
	} else {
		log.Println("Connected to remote repository: " + repoName)
	}

	// Use the default registry credentials
	username := c.registryUsername
	password := c.registryPassword
	log.Println("Using registry credentials: " + username + ":" + password)

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: auth.StaticCredential(registryDomain, auth.Credential{
			Username: username,
			Password: password,
		}),
	}

	// 5. Copy from the file store to the remote repository
	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	if err != nil {
		fmt.Println("Error pushing files from: " + dir + " to " + repo.Reference.Repository + ":" + tag)
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		log.Println("Pushed files from: " + dir + " to " + repo.Reference.Repository + ":" + tag)
	}
}
