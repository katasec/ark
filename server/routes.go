package server

import (
	"net/http"

	"github.com/katasec/ark/server/handlers"
)

func (s *Server) initaliseRoutes() {

	// Setup file server
	fs := http.FileServer(http.Dir("server/assets"))
	s.router.Handle("/assets/*", http.StripPrefix("/assets", fs))

	// Register route handlers for routes
	s.router.Get("/", handlers.HomeHander(s))
	s.router.Post("/azure/cloudspace", handlers.PostCloudspace(s))
	s.router.Delete("/azure/cloudspace", handlers.DeleteCloudspace(s))
	s.router.Post("/azure/vm", handlers.PostVm(s))

}
