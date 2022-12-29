package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

var (
	chiDispatcher = chi.NewRouter()
)

type chiRouter struct{}

func NewChiRouter() ArkRouter {
	chiDispatcher.Use(middleware.Logger)
	return &chiRouter{}
}

// GET implements ArkRouter
func (*chiRouter) GET(path string, f http.HandlerFunc) {
	chiDispatcher.Get(path, f)
}

// POST implements ArkRouter
func (*chiRouter) POST(path string, f http.HandlerFunc) {
	chiDispatcher.Post(path, f)
}

// SERVE implements ArkRouter
func (*chiRouter) LISTEN(port string) {
	log.Printf("Chi router running on port %s \n", port)
	log.Fatal(http.ListenAndServe(port, chiDispatcher))
}
