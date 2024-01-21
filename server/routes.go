package server

import (
	"net/http"
)

func (s *Server) initaliseRoutes() {

	// Setup file server
	fs := http.FileServer(http.Dir("server/assets"))
	s.router.Handle("/assets/*", http.StripPrefix("/assets", fs))

	// Register route handlers for routes
	s.router.Get("/", s.HomeHander())

	// Cloudspace routes
	s.router.Post("/azure/cloudspace", s.PostCloudspace())
	s.router.Delete("/azure/cloudspace", s.DeleteCloudspace())

	// Hello routes
	s.router.Post("/azure/hello", s.PostHello())

	s.router.Post("/azure/vm", s.PostVm())

}
