package router

import "net/http"

type ArkRouter interface {
	GET(uri string, f http.HandlerFunc)
	POST(uri string, f http.HandlerFunc)
	LISTEN(port string)
}
