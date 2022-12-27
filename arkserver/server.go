package arkserver

import (
	"fmt"
	"net/http"

	"github.com/katasec/ark/arkrouter"
)

func Start() {
	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })
	// http.ListenAndServe(":3000", r)

	r := arkrouter.NewChiRouter()

	r.GET("/", home)

	r.SERVE(":8080")
}

func home(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello World!")
}
