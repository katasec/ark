package web

import (
	"net/http"

	"github.com/katasec/ark/web/handlers"
)

func (s *Server) initialiseRoutes() {

	s.mux.Handle("/assets/", run(handlers.FileHandler("/assets")))

	s.mux.Handle("/api/arkweb", run(handlers.HomeHandler()))
	s.mux.Handle("/api/web", run(handlers.HomeHandler()))

	s.mux.Handle("/api/arkweb/", run(handlers.FileHandler("/api/arkweb/")))
	s.mux.Handle("/api/web/", run(handlers.FileHandler("/api/arkweb/")))
	// Serve "assets" folder from root
	s.mux.Handle("/", run(handlers.HomeHandler()))

}

func run(h http.HandlerFunc) http.HandlerFunc {
	return MultipleMiddleware(
		h,
		handlers.LogHandlerFunc,
	)

}
func getMiddlewareSequence() http.HandlerFunc {
	return MultipleMiddleware(
		handlers.EntryHandler,
		handlers.LogHandlerFunc,
		handlers.FileHandlerFunc("/"),
	)
}
