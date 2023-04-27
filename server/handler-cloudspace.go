package server

import (
	"fmt"
	"net/http"

	"github.com/katasec/ark/requests"
	"gopkg.in/yaml.v2"
)

func (s *Server) postCloudspace() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set default acs name
		acsRequest := requests.CreateAzureCloudspaceRequest{}
		if acsRequest.Name == "" {
			acsRequest.Name = "default"
		}

		// Decode request body into acsRequest
		err := yaml.NewDecoder(r.Body).Decode(&acsRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Type", "application/x-yaml")
		//s.acsrepo.GetCloudSpace(acsRequest.Name)

		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, acsRequest.ToYamlAzureCloudpace())

		s.msg.Send("azurecloudspace", acsRequest.ToYamlAzureCloudpace())
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

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, request.ToYamlAzureCloudpace())

		s.msg.Send("DeleteAzureCloudspaceRequest", request.ToYamlAzureCloudpace())
	})
}
