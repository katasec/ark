package push

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func DoPush(url string, tag string) {
	tmpdir := CloneRemote(url, tag)
	err := os.Chdir(tmpdir)
	if err != nil {
		log.Println("Could not change dir to tmpdir: " + tmpdir)
	} else {
		log.Println("Changed dir to tmpdir: " + tmpdir)
	}

	err = tarAndGzip(tmpdir)
	if err != nil {
		fmt.Println("Error creating tar archive:", err)
		os.Exit(1)
	}
}

func CloneRemote(url string, tag string) string {

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
	log.Println("Repo Dir: " + tmpdir)

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

func Fprintln(w io.Writer, message string) {
	t := time.Now()
	message = fmt.Sprint(t.Format("2006/01/02 15:04:05") + " " + message)
	fmt.Fprintln(w, message)
}

// Function to process a file or directory during traversal
func processFile(tarWriter *tar.Writer, path string, info os.FileInfo) error {
	if info.IsDir() && path == ".git" {
		return filepath.SkipDir // Skip the ".git" directory and its contents
	}

	if info.Mode().IsRegular() {

		// Sip git files
		if !strings.Contains(path, ".git/") {
			// Create a new tar header, using FileInfo data
			header, err := tar.FileInfoHeader(info, path)
			if err != nil {
				return err
			}

			// Update the header's name to correctly reflect the desired destination when untaring
			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			// Open files for taring
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy file data into tar writer
			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

	}

	return nil
}

// Function to create a tar archive of a directory and compress it with gzip
func tarAndGzip(sourceDir string) error {

	targetFile := sourceDir + ".tar.gz"
	tarFile, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the target file itself
		if path == targetFile {
			return nil
		}

		return processFile(tarWriter, path, info) // Call the separate function
	})

	return err
}

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
	if err != nil {
		panic(err)
	} else {
		fmt.Println(repo.Reference.Repository + ":" + tagx + " copied to " + localpath)
	}

}
