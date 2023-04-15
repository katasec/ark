package server

import "net/http"

func (s *Server) initaliseRoutes() {
	// Register route handlers for routes
	s.router.Get("/", s.homeHander())
	s.router.Post("/azure/cloudspace", s.postCloudspace())
	s.router.Post("/azure/vm", s.postVm())

	fileServer := http.FileServer(http.Dir("./images/"))
	s.router.Handle("/images/", http.StripPrefix("/images", fileServer))

}
