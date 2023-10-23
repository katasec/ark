package server

import (
	"fmt"
	"log"
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

		// Get cloudspace from DB
		acs, err := s.acsrepo.GetCloudSpace(acsRequest.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error getting cloudspace: %s,", err)
			return
		} else {
			fmt.Println("Cloudspace found in DB")
		}

		if acs.Name == "" {
			acs.Name = acsRequest.Name
		}

		// Add environments from reuqest
		for _, env := range acsRequest.Environments {
			acs.AddSpoke(env)
		}

		fmt.Println("The hub's name:", acs.Hub.Name)

		// Send request to queue
		err = s.qClient.Send("azurecloudspace", acs.ToYaml())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Error: %s,", err)
			return
		} else {
			log.Println("Message sent successfully !")
		}

		// Save cloudspace to DB
		//s.acsrepo.CreateCloudSpace(&acs)
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
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

		s.qClient.Send("DeleteAzureCloudspaceRequest", request.ToYamlAzureCloudpace())

	})
}
