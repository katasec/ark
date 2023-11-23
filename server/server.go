package server

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
	"github.com/katasec/ark/repositories"

	_ "github.com/lib/pq" // Import the pq driver
)

// Server struct models the ark server and its dependencies
type Server struct {
	router  *chi.Mux
	config  *config.Config
	cmdQ    messaging.Messenger
	respQ   messaging.Messenger
	db      *sql.DB
	Acsrepo *repositories.AzureCloudSpaceRepository
}

func NewServer() *Server {

	// Read from local config  file
	cfg := config.ReadConfig()

	// Setup Router
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Logger)

	//Setup DB Config
	db, err := getDbConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize Repositories
	acsrepo := repositories.NewAzureCloudSpaceRepository(db)

	// Return server with local config
	return &Server{
		config:  cfg,
		cmdQ:    messaging.NewRabbitMqMessenger(cfg.MqConnStr, cfg.CmdQ),
		respQ:   messaging.NewRabbitMqMessenger(cfg.MqConnStr, cfg.RespQ),
		router:  chiRouter,
		Acsrepo: acsrepo,
	}
}

func (s *Server) Start() {

	// Initialize Routes
	s.initaliseRoutes()

	// Start Listening
	log.Println("Server started on port " + s.config.ApiServer.Port + "...")
	log.Fatal(http.ListenAndServe(":"+s.config.ApiServer.Port, s.router))

}

func (s *Server) getDbConnection() (*sql.DB, error) {

	db, err := sql.Open(s.config.DbConfig.DriverName, s.config.DbConfig.DataSourceName)
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	} else {
		log.Println("Database ping successful!")
	}

	return db, err
}

func getDbConnection() (*sql.DB, error) {

	config := config.ReadConfig()

	db, err := sql.Open(config.DbConfig.DriverName, config.DbConfig.DataSourceName)
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	} else {
		log.Println("Database ping successful!")
	}

	return db, err
}

// Funcstions to impolement Server interface

func (s *Server) GetCommandQ() messaging.Messenger {
	return s.cmdQ
}

func (s *Server) GetRouter() *chi.Mux {
	return s.router
}

func (s *Server) GetConfig() *config.Config {
	return s.config
}

func (s *Server) GetDb() *sql.DB {
	return s.db
}

func (s Server) GetAcsrepo() *repositories.AzureCloudSpaceRepository {
	return s.Acsrepo
}
