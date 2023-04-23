package server

import (
	"fmt"
	"net/http"

	"github.com/katasec/ark/requests"
	"gopkg.in/yaml.v2"
)

func (s *Server) postCloudspace() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request := requests.CreateAzureCloudspaceRequest{}

		err := yaml.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, request.ToJsonAzureCloudpace())

		s.msg.Send("azurecloudspace", request.ToJsonAzureCloudpace())
	})
}

func (s *Server) deleteCloudspace() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request := requests.DeleteAzureCloudspaceRequest{
			Name: "default",
		}

		err := yaml.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, request.ToJsonAzureCloudpace())

		s.msg.Send("DeleteAzureCloudspaceRequest", request.ToJsonAzureCloudpace())
	})
}
