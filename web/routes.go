package web

import (
	"net/http"

	"github.com/katasec/ark/web/handlers"
)

func (s *Server) initialiseRoutes() {

	// Define the sequence of middleware used to process all requests
	//middlewareSequence := getMiddlewareSequence()

	// Use the file system to serve static files
	//s.mux.Handle("/", middlewareSequence)

	//fs := http.FileServer(http.FS(handlers.GetStaticAssets()))

	//patterns := []string{"/", "/assets"}
	// for _, pattern := range patterns {
	// 	stripped := http.StripPrefix(pattern, handlers.LogHandler(handlers.FileHandler()))
	// 	s.mux.Handle(pattern, stripped)
	// }

	s.mux.Handle("/assets/", run(handlers.FileHandler("/assets")))
	s.mux.Handle("/api/arkweb", run(handlers.HomeHandler()))
	s.mux.Handle("/api/arkweb/", run(handlers.FileHandler("/api/arkweb/")))

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
