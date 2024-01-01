package arkimage

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// isValidGitUrl checks if a URL is a valid Git URL
func isValidGitUrl(repoUrl string) bool {

	// Currently just accepting Github URLs and is very simplistic

	// Parse the URL using the standard library
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		return false
	}

	// Check for valid Git transport schemes
	if parsedURL.Scheme != "https" && parsedURL.Scheme != "git" && parsedURL.Scheme != "ssh" {
		return false
	}

	// Check for common Git hosts and extensions
	if !strings.HasPrefix(parsedURL.Host, "github.com") {
		return false
	}

	// Additional checks for SSH URLs
	if parsedURL.Scheme == "ssh" {
		if parsedURL.User == nil || parsedURL.Path == "" {
			return false
		}
	}

	_, valid := hasValidResourceName(repoUrl)
	return valid
}

func hasValidResourceName(repoUrl string) (resourceName string, valid bool) {

	// Split the URL into parts by "/" and get the last fragment.
	// For e.g  https://github.com/katasec/ark-resource-azurecloudspace will return "ark-resource-azurecloudspace"
	parts := strings.Split(repoUrl, "/")
	repoName := parts[len(parts)-1]

	// Check it starts with "ark-resource-"
	if !strings.HasPrefix(repoName, "ark-resource-") {
		return "", false
	}

	resourceName = strings.TrimPrefix(repoName, "ark-resource-")
	return resourceName, true
}

// cloneRemote clones a remote repo into a temp dir
func cloneRemote(url string, tag string) string {
	// Create a temp dir
	tmpdir := createTempDir()

	// Clone Repo
	repo, err := git.PlainClone(tmpdir, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		log.Println("Cloning error:" + err.Error())
	} else {
		log.Println("Cloned: " + url)
	}

	// Get worktree
	w, err := repo.Worktree()
	if err != nil {
		log.Println("Worktree error:" + err.Error())
		os.Exit(1)
	}

	// Checkout tag
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + tag),
		Force:  true,
	})
	if err != nil {
		log.Printf("Error checking out tag: %v. Message: %v\n", tag, err.Error())
		os.Exit(1)
	} else {
		log.Println("Checked out tag: " + tag)
	}

	return tmpdir
}
