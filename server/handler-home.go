package server

import (
	"fmt"
	"net/http"
)

var (
	homeHtml = `
		<!DOCTYPE html>
		<html>
			<head>
				<title>Home</title>
				</head>
				<body>
					<img src="/images/ark.svg" alt="Ark Logo" />
				</body>
		</html>
					`
)

func (s *Server) homeHander() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, homeHtml)
	})
}
