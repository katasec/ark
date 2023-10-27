package server

import (
	"fmt"
	"net/http"
)

var (
	link     = "https://katasecid.b2clogin.com/katasecid.onmicrosoft.com/oauth2/v2.0/authorize?p=B2C_1_SignUpIn&client_id=6fa18972-6f2d-4bca-a4a6-0ebd58f77bf6&nonce=defaultNonce&redirect_uri=https%3A%2F%2Fjwt.ms&scope=openid&response_type=id_token&prompt=login"
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
						<img src="/images/ark.svg" alt="Ark Logo" width="10%"/ >
					</div>

					<div id="desription"  class="style1">
						<br/>
						Ark is a cli to simplify management of your cloud resources using security best practices. Please check the <a href="https://github.com/katasec/ark">repo</a> for more details.
						Click <a href="https://katasecid.b2clogin.com/katasecid.onmicrosoft.com/oauth2/v2.0/authorize?p=B2C_1_SignUpIn&client_id=6fa18972-6f2d-4bca-a4a6-0ebd58f77bf6&nonce=defaultNonce&redirect_uri=https%3A%2F%2Fjwt.ms&scope=openid&response_type=id_token&prompt=login">here</a> to sign up!
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
