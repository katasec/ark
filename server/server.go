package server

import (
	"fmt"
	"net/http"

	"github.com/katasec/ark/config"
	"github.com/katasec/ark/router"
)

// Server struct models the ark server and its dependencies
type Server struct {
	router *router.ArkRouter
	config *config.Config
}

func NewServer() *Server {

	router := router.NewChiRouter()
	cfg := config.ReadConfig()

	return &Server{
		router: &router,
		config: cfg,
	}
}

func (*Server) Start() {

	r := router.NewChiRouter()

	r.GET("/", home)

	r.LISTEN(ListenPort)
}

func home(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello World!")
}
