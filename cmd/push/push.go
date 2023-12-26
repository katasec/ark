package push

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func DoPush(url string, tag string) {

	// Clone repo into a temp dir
	repoDir := cloneRemote(url, tag)
	os.Chdir(repoDir)

	// Push the cloned directory to the registry
	pushArchiveToRegistry(repoDir)
}

// cloneRemote clones a remote repo into a temp dir
func cloneRemote(url string, tag string) string {
	// Create a temp dir
	tmpdirBase := filepath.Join(os.TempDir(), "ark")
	err := os.Mkdir(tmpdirBase, os.FileMode(0777))
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		fmt.Println("could not create tmpdirBase, exitting." + tmpdirBase)
		fmt.Println(err.Error())
		os.Exit(1)
	}
	tmpdir, _ := os.MkdirTemp(tmpdirBase, "ark-remote")
	log.Println("Cloning: " + url)

	// Clone Repo
	_, err = git.PlainClone(tmpdir, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		log.Println("Cloning error:" + err.Error())
	} else {
		log.Println("Done.")
	}

	return tmpdir
}

// tarAndGzip Creates a tar archive of a directory and compress it with gzip
func tarAndGzip(sourceDir string) *os.File {
	// Create the tar file
	targetFile := sourceDir + ".tar.gz"
	tarFile, err := os.Create(targetFile)
	if err != nil {
		fmt.Println("Error creating tar file:", err)
		os.Exit(1)
	}
	defer tarFile.Close()

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Recursively walk through the directory and add files to the tarWriter
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error walking path:", err)
			os.Exit(1)
		}

		// Skip the target file itself
		if path == targetFile {
			return nil
		}
		return addFileToTar(tarWriter, path, info) // Call the separate function
	})

	if err != nil {
		fmt.Println("Error creating tar archive:", err)
		os.Exit(1)
	}

	log.Println("The created tar file was:" + tarFile.Name())
	return tarFile
}

// addFileToTar Adds a file to a tar archive
func addFileToTar(tarWriter *tar.Writer, path string, info os.FileInfo) error {
	// Skip directories, files in the .git directory and the .gitignore file
	if strings.Contains(path, ".git/") || strings.Contains(path, ".gitignore") || !info.Mode().IsRegular() {
		return nil
	}

	// Create tar header for file
	header, err := tar.FileInfoHeader(info, path)
	if err != nil {
		return err
	}

	// Write header to tarWriter
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// Open file for taring
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy file data into tarWriter
	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}
	return nil
}

// pushArchiveToRegistry pushes a tar.gz to a registry
func pushArchiveToRegistry(tmpdirBase string) {

	// Pusing a file to a registry requires the creation of an ORA
	// local file store and a remote repository. The file store is a local directory
	// where the file is stored. The remote repository is the registry where the file is pushed to.

	// 0. Create a file store
	fs, err := file.New(tmpdirBase)
	log.Println("Creating file store: " + tmpdirBase)
	if err != nil {
		fmt.Println("Error creating file store:", err)
		os.Exit(1)
	}
	defer fs.Close()
	ctx := context.Background()

	// 1. Add files to the file store
	mediaType := "application/vnd.test.file"
	fileNames := listFilesRecursively(tmpdirBase)
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))

	for _, name := range fileNames {
		// name is the absolute path to the file and we want to store it relative
		// to the file store. As such we trim the tmpdirBase prefix.
		relativePath := strings.TrimPrefix(name, tmpdirBase+"/")
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
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1_RC4, artifactType, oras.PackManifestOptions{
		Layers: fileDescriptors,
	})
	if err != nil {
		fmt.Println("Error packing manifest:", err)
		os.Exit(1)
	}

	tag := "latest"
	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		fmt.Println("Error tagging manifest:", err)
		os.Exit(1)
	}

	// 3. Connect to a remote repository

	// Get Remote registry details
	ics := "ghcr.io/katasec/cloudspace"
	ref := strings.Split(ics, ":")[0]
	//tagx := strings.Split(ics, ":")[1]
	tagx := "latest"
	registryDomain := strings.Split(ics, "/")[0]

	// Connect to the remote repository
	repo, err := remote.NewRepository(ref)
	if err != nil {
		panic(err)
	} else {
		log.Println("Connected to remote repository: " + ref)
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

	// 4. Copy from the file store to the remote repository
	src := fs
	dst := repo
	_, err = oras.Copy(ctx, src, tagx, dst, tagx, oras.DefaultCopyOptions)
	if err != nil {
		fmt.Println("Error pushing files from: " + tmpdirBase + " to " + repo.Reference.Repository + ":" + tagx)
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		log.Println("Pushed files from: " + tmpdirBase + " to " + repo.Reference.Repository + ":" + tagx)
	}
}

func listFilesRecursively(dirPath string) []string {
	var files []string

	fmt.Println("Recusively listing files in: " + dirPath)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// Check error walking path
		if err != nil {
			log.Println("Error walking path:", err)
			os.Exit(1)
		}

		// Add file to files if it is not a directory
		if !strings.Contains(path, ".git/") && !strings.Contains(path, ".gitignore") && !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Println("Error walking path:", err)
		os.Exit(1)
	}

	return files
}
