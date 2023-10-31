package handlers

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func HomeHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Open the HTML file.
		htmlBytes, err := fs.ReadFile(assets, "assets/index.html")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		html := string(htmlBytes)

		// Create a template with a variable placeholder
		tmpl := template.Must(template.New("home").Parse(html))
		data := map[string]string{
			"prefix": storagePrefix,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Fatal(err)
		}
	}
}
