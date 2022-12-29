package arkserver

import (
	"fmt"
	"net/http"

	"github.com/katasec/ark/arkrouter"
	"github.com/katasec/ark/config"
)

// Server struct models the wire serer and its dependencies
type Server struct {
	// config       config.Config
	router *arkrouter.ArkRouter
	config *config.Config
}

func NewServer() *Server {

	router := arkrouter.NewChiRouter()
	cfg := config.ReadConfig()

	return &Server{
		router: &router,
		config: cfg,
	}
}

func Start() {

	r := arkrouter.NewChiRouter()

	r.GET("/", home)

	r.SERVE(":8080")
}

func home(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello World!")
}
