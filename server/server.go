package server

import (
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/database"
	resources "github.com/katasec/ark/resources"
	"github.com/katasec/ark/router"
)

// Server struct models the ark server and its dependencies
type Server struct {
	router router.ArkRouter
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

	repo := database.NewJsonRepository()
	repo.AddVm(resources.Vm{
		ProjectName: "vmproject",
		Name:        "vm01",
		Tags: map[string]string{
			"role": "vm",
		},
	})

	repo.AddVm(resources.Vm{
		ProjectName: "myvms",
		Name:        "vm02",
	})

	repo.SaveVms()
	return

	// Select Router type (For e.g. Chi vs. Gorilla mux)
	s.router = router.NewChiRouter()

	// Initialize Routes
	s.initaliseRoutes()

	// Start Listening
	s.router.LISTEN(ListenPort)
}

func (s *Server) DbStuff() {
	database.SomeStuff()
}
