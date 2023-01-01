package server

import "github.com/katasec/ark/database"

var (
	ListenPort = ":8080"
	db         = database.NewJsonRepository()
)
