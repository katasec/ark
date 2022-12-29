package server

import "net/http"

func (s *Server) initaliseRoutes() {

	// Register route handlers for routes
	s.router.GET("/", s.homeHander())

	s.azureGet("home", s.homeHander())
}

// azureRoute registers an Azure Route with the router
func (s *Server) azureGet(uri string, f http.HandlerFunc) {
	// Register route handlers for routes
	s.router.GET("/azure/"+uri, s.homeHander())
}
