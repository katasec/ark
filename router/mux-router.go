package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	r = mux.NewRouter()
)

type muxRouter struct{}

func NewMuxRouter() ArkRouter {
	return &muxRouter{}
}

// GET implements ArkRouter
func (*muxRouter) GET(path string, f http.HandlerFunc) {
	r.HandleFunc(path, f).Methods("GET")
}

// POST implements ArkRouter
func (*muxRouter) POST(path string, f http.HandlerFunc) {
	r.HandleFunc(path, f).Methods("POST")
}

// SERVER implements ArkRouter
func (*muxRouter) LISTEN(port string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Println("Starting Mux server")
	log.Fatal(http.ListenAndServe(port, loggedRouter))
}
