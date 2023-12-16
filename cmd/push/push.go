package push

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

func DoPush(url string, tag string) {
	CloneRemote(url, tag)
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
