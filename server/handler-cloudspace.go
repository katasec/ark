package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	resources "github.com/katasec/ark/resources/v0"
)

func (s *Server) postCloudspace() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cloudspace := resources.CloudSpace{}

		err := json.NewDecoder(r.Body).Decode(&cloudspace)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Cloud Space: %+v", cloudspace)

		db.AddCloudSpace(cloudspace)
		db.SaveCloudSpaces()
	})
}
