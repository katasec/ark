package handlers

import (
	"fmt"
	"net/http"
)

func LogHandlerFunc(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logger: " + r.Host + " " + r.Method + " " + r.URL.Path + " " + r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func LogHandler(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logger: " + r.Host + " " + r.Method + " " + r.URL.Path + " " + r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}
