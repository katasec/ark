package arkimage

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

type ArkImage struct {
	ctx context.Context
}

func NewArkImage() *ArkImage {
	return &ArkImage{
		ctx: context.Background(),
	}
}

// Push Pushes code from a git repo to a registry
func (c *ArkImage) Push(gitUrl string, tag string) {
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

func (c *ArkImage) GetLocalPath(image string) string {
	// Get home directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	fqImage := getFqImage(image)
	resourceName := getResourceName(fqImage)
	version := getVersion(fqImage)

	localpath := path.Join(homedir, ".ark", "registry", arkConfig.ArkRegistry.Domain, arkConfig.ArkRegistry.RepoName, resourceName, version)

	return localpath
}

// Pull Pulls code from a registry to a local directory
func (c *ArkImage) Pull(image string) {

	//ghcr.io/katasec/ark-resource-hello:v0.0.1 or ark-resource-hello:v0.0.1

	fqImage := getFqImage(image)

	// Create local path under ~/.ark/registry/registrydomain/reponame/resourcename/version/
	localpath := c.GetLocalPath(image)
	//log.Println("The localpath is:" + localpath)

	// Create a file store in the local path
	fs, err := file.New(localpath)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// Delete files in file store if any
	//log.Println("Deleting files in file store if any before download: " + localpath)
	deleteDirectoryContents(localpath)

	// Connect to a remote repository
	ref := strings.Split(fqImage, ":")[0]
	repo, err := remote.NewRepository(ref)
	if err != nil {
		panic(err)
	}
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: auth.StaticCredential(arkConfig.ArkRegistry.Domain, auth.Credential{
			Username: arkConfig.ArkRegistry.Username,
			Password: arkConfig.ArkRegistry.Password,
		}),
	}

	// Pull the image to the local file store
	tagx := strings.Split(image, ":")[1]
	_, err = oras.Copy(c.ctx, repo, tagx, fs, tagx, oras.DefaultCopyOptions)
	if err != nil {
		log.Println("Error pulling " + repo.Reference.Repository + ":" + tagx)
		log.Println(err.Error())
	} else {
		//fmt.Println(repo.Reference.Repository + ":" + tagx + " copied to " + localpath)
	}

	// delete file if it exists
	log.Println("Deleting configdata.json if it exists")
	os.Remove(localpath + "/configdata.json")

}

func getFqImage(image string) string {
	if strings.Contains(image, "/") {
		return image
	} else {
		return arkConfig.ArkRegistry.Domain + "/" + arkConfig.ArkRegistry.RepoName + "/" + image
	}
}

func getResourceName(fqImage string) string {
	resourceName := strings.Split(fqImage, "/")[len(strings.Split(fqImage, "/"))-1]
	resourceName = strings.Split(resourceName, ":")[0]
	return resourceName
}

func getVersion(fqImage string) string {
	resourceName := strings.Split(fqImage, "/")[len(strings.Split(fqImage, "/"))-1]
	resourceName = strings.Split(resourceName, ":")[1]
	return resourceName
}

// pushToRegistry Pushes files from a directory to a registry
func (c *ArkImage) pushToRegistry(dir string, tag string, gitUrl string) {

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
	log.Println("Resource name: " + resourceName)
	repoName := arkConfig.ArkRegistry.Domain + "/" + arkConfig.ArkRegistry.RepoName + "/" + resourceName
	repoName = strings.Replace(repoName, "//", "/", -1)
	registryDomain := arkConfig.ArkRegistry.Domain
	log.Println("Registry domain: " + registryDomain)

	// 4. Connect to the remote repository
	log.Println("Connecting to: " + repoName)
	repo, err := remote.NewRepository(repoName)
	if err != nil {
		panic(err)
	} else {
		log.Println("Connected to remote repository: " + repoName)
	}

	// Use the default registry credentials
	username := arkConfig.ArkRegistry.Username
	password := arkConfig.ArkRegistry.Password

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
