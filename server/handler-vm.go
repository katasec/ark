package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	resources "github.com/katasec/ark/resources"
)

func (s *Server) PostVm() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vm := resources.Vm{}

		err := json.NewDecoder(r.Body).Decode(&vm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Vm: %+v", vm)

		// db.AddVm(vm) // Add to memory
		// db.SaveVms() // Write to disk

	})
}
