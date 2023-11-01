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
		port: misc.IfEmpty(os.Getenv("ARK_WEB_PORT"), "80"),
	}

	// Register handlers for routes
	s.initialiseRoutes()

	return s
}

func (s *Server) Start() {
	// Start the server
	listAddr := "0.0.0.0"
	log.Println("Listening on: http://" + listAddr + ":" + s.port)
	log.Fatal(http.ListenAndServe(listAddr+":"+s.port, s.mux))
}
