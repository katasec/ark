package arkimage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
func (c *ArkImage) Push(gitUrl string, tag string, imageType string) {
	fmt.Println("Pushing crate")

	if isValidGitUrl(gitUrl) {
		// Clone repo into a temp dir
		repoDir := cloneRemote(gitUrl, tag)
		os.Chdir(repoDir)

		// Push the cloned directory to the registry
		c.pushToRegistry(repoDir, tag, gitUrl, imageType)
	} else {
		fmt.Println("Not a git URL: " + gitUrl)
	}

}

// Pull Pulls code from a registry to a local directory
func (c *ArkImage) Pull(image string) string {

	imageType := "" // For e.g. pulumi or terraform

	// Create local path under ~/.ark/registry/...
	localpath := c.GetLocalPath(image)

	// Create a file store in the local path
	fs, err := file.New(localpath)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// Delete files in file store if any
	deleteDirectoryContents(localpath)

	// Get fully qualified image name (with registry domain)
	fqImage := getFqImage(image)

	// Create repository struct to connect to registry
	ref := strings.Split(fqImage, ":")[0]
	repo, err := remote.NewRepository(ref)
	if err != nil {
		panic(err)
	}

	// Auth creds for registry
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: auth.StaticCredential(arkConfig.ArkRegistry.Domain, auth.Credential{
			Username: arkConfig.ArkRegistry.Username,
			Password: arkConfig.ArkRegistry.Password,
		}),
	}

	// Pull/Copy image from repo to the local file store
	tagx := strings.Split(image, ":")[1]
	descriptor, err := oras.Copy(c.ctx, repo, tagx, fs, tagx, oras.DefaultCopyOptions)
	if err != nil {
		log.Println("Error pulling " + repo.Reference.Repository + ":" + tagx)
		log.Println(err.Error())
		os.Exit(1)
	} else {
		// Get manifest reader after image pull
		manifestReader, err := fs.Fetch(c.ctx, descriptor)
		if err != nil {
			panic(err)
		}

		// Get manifest data
		data, err := io.ReadAll(manifestReader)
		if err != nil {
			panic(err)
		}

		// Unmarshal the manifest to access annotations
		var manifest Manifest
		err = json.Unmarshal(data, &manifest)
		if err != nil {
			panic(err) // Handle error appropriately
		}

		imageType = manifest.Annotations.ImageType
	}

	// delete file if it exists
	os.Remove(localpath + "/configdata.json")

	return imageType
}

// pushToRegistry Pushes files from a directory to a registry
func (c *ArkImage) pushToRegistry(dir string, tag string, gitUrl string, imageType string) {

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
		ManifestAnnotations: map[string]string{
			"org.opencontainers.image.type": imageType,
		},
	}
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_0, artifactType, opts)
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
	log.Println("Reponame:" + repoName)
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

// getFqImage Returns the fully qualified image name that include regis domain and repo name
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
