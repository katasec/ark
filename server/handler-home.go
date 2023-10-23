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
					<style>
					.style1
					{
						font-family:    Arial, Helvetica, sans-serif;
						font-size:      15px;
						text-align : center;
					}
				</style>					
				</head>
				<body>
					<div id="logo" class="style1">
						<img src="/images/ark.svg" alt="Ark Logo" width="100"/>
					</div>

					<div id="desription"  class="style1">
						<br/>
						Ark is a cli to simplify management of your cloud resources using security best practices. Please check the <a href="https://github.com/katasec/ark">repo</a> for more details.
					</div>					
				</body>
		</html>
					`
)

func (s *Server) homeHander() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, homeHtml)
	})
}
