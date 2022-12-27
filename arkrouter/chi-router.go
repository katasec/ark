package arkrouter

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
func (*chiRouter) GET(path string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(path, f)
}

// POST implements ArkRouter
func (*chiRouter) POST(path string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(path, f)
}

// SERVE implements ArkRouter
func (*chiRouter) SERVE(port string) {
	log.Printf("Chi router running on port %s \n", port)
	http.ListenAndServe(port, chiDispatcher)
}
