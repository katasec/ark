package web

import (
	"log"
	"net/http"
	"os"

	"github.com/katasec/utils/misc"
)

type Server struct {
	port string
	mux  *http.ServeMux
}

func NewServer() *Server {

	// Create Server struct with new mux
	s := &Server{
		mux:  http.NewServeMux(),
		port: misc.IfEmpty(os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT"), "8080"),
	}

	// Register handlers for routes
	s.initialiseRoutes()

	return s
}

func (s *Server) Start() {
	// Start the server
	log.Println("Listening on:" + s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, s.mux))
}
