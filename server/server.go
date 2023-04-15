package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/database"
)

// Server struct models the ark server and its dependencies
type Server struct {
	router *chi.Mux
	config *config.Config
}

func NewServer() *Server {

	// Read from local config  file
	cfg := config.ReadConfig()

	// Return server with local config
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() {

	fmt.Println("Starting server")

	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Logger)

	s.router = chiRouter

	// Initialize Routes
	s.initaliseRoutes()

	// Start Listening
	log.Fatal(http.ListenAndServe(":"+s.config.ApiServer.Port, s.router))

}

func (s *Server) DbStuff() {
	database.SomeStuff()
}
