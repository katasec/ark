package handlers

import (
	"fmt"
	"net/http"
)

// EntryHandler ...
func EntryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I am in Entry Handler!")
}
