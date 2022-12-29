package server

import (
	"fmt"
	"net/http"
)

func (s *Server) homeHander() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
}
