package server

import (
	"fmt"
	"net/http"

	"github.com/katasec/ark/requests"
	"gopkg.in/yaml.v2"
)

func (s *Server) postCloudspace() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request := requests.AzureCloudspaceRequest{}

		err := yaml.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, env := range request.Environments {
			fmt.Println("Env:" + env)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, request.ToAzureCloudpace())

		s.msg.Send("azurecloudspace", request.ToAzureCloudpace())
	})
}
