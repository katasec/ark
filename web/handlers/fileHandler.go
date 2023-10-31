package handlers

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

func FileHandlerFunc(prefix string) func(http.HandlerFunc) http.HandlerFunc {

	return func(http.HandlerFunc) http.HandlerFunc {

		// Use the file system to serve static files
		assets := GetStaticAssets()
		fs := http.FileServer(http.FS(assets))

		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("In File Handler:" + r.URL.RawPath)
			strippedRequest := r.Clone(r.Context())
			strippedRequest.URL.Path = strings.TrimPrefix(strippedRequest.URL.Path, prefix)

			fs.ServeHTTP(w, r)
		}
	}
}
func FileHandlerMiddleWare(next http.HandlerFunc) http.HandlerFunc {

	// Use the file system to serve static files
	assets := GetStaticAssets()
	fs := http.FileServer(http.FS(assets))

	return func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}
}
func FileHandler(prefix ...string) http.HandlerFunc {

	var stripPrefix string
	if len(prefix) > 0 {
		stripPrefix = prefix[0]
	} else {
		stripPrefix = ""
	}

	// Use the file system to serve static files
	assets := GetStaticAssets()
	fs := http.FileServer(http.FS(assets))

	if stripPrefix != "" {
		fmt.Println("Returning stripped !")
		return func(w http.ResponseWriter, r *http.Request) {
			strippedRequest := r.Clone(r.Context())
			strippedRequest.URL.Path = strings.TrimPrefix(strippedRequest.URL.Path, stripPrefix)
			fmt.Println("The stripped path was:" + strippedRequest.URL.Path)
			fs.ServeHTTP(w, strippedRequest)
		}

	} else {
		fmt.Println("Returning normal !")
		return func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}
	}

}

func GetStaticAssets() fs.FS {
	fs, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}

	return fs
}
