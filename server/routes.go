package server

func (s *Server) initaliseRoutes() {
	// Register route handlers for routes
	s.router.GET("/", s.homeHander())
	s.router.POST("/azure/cloudspace", s.postCloudspace())
}
