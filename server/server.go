package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/database"
	"github.com/katasec/ark/messaging"
)

// Server struct models the ark server and its dependencies
type Server struct {
	router *chi.Mux
	config *config.Config
	msg    messaging.Messenger
}

func NewServer() *Server {

	// Read from local config  file
	cfg := config.ReadConfig()

	// Setup Router
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Logger)

	// Setup messaging client
	connString := cfg.AzureConfig.MqConfig.MqConnectionString
	queueName := cfg.AzureConfig.MqConfig.MqName
	msg := messaging.NewAsbMessenger(connString, queueName)

	// Return server with local config
	return &Server{
		config: cfg,
		msg:    msg,
		router: chiRouter,
	}
}

func (s *Server) Start() {

	// Initialize Routes
	s.initaliseRoutes()

	// Start Listening
	log.Println("Server started on port " + s.config.ApiServer.Port + "...")
	log.Fatal(http.ListenAndServe(":"+s.config.ApiServer.Port, s.router))

}

func (s *Server) DbStuff() {
	database.SomeStuff()
}
