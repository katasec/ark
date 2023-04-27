package server

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
	"github.com/katasec/ark/repositories"
)

// Server struct models the ark server and its dependencies
type Server struct {
	router *chi.Mux
	config *config.Config
	msg    messaging.Messenger
	db     *sql.DB

	acsrepo *repositories.AzureCloudSpaceRepository
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

	// Setup Database
	dbDir := cfg.GetDbDir()
	dbFile := fmt.Sprintf("%s/ark.db", dbDir)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Database opened successfully")
	}

	// Initialize Repositories
	acsrepo := repositories.NewAzureCloudSpaceRepository(db)

	// Return server with local config
	return &Server{
		config:  cfg,
		msg:     msg,
		router:  chiRouter,
		db:      db,
		acsrepo: acsrepo,
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
	// database.SomeStuff()
}
