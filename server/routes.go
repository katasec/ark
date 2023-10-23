package server

import "net/http"

func (s *Server) initaliseRoutes() {

	// Setup file server
	fs := http.FileServer(http.Dir("images"))
	s.router.Handle("/images/*", http.StripPrefix("/images/", fs))

	// Register route handlers for routes
	s.router.Get("/", s.homeHander())
	s.router.Post("/azure/cloudspace", s.postCloudspace())
	s.router.Delete("/azure/cloudspace", s.deleteCloudspace())
	s.router.Post("/azure/vm", s.postVm())
}
