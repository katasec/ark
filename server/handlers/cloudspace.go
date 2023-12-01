package handlers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	logx "github.com/katasec/ark/log"

	"github.com/katasec/ark/arkserver"
	"github.com/katasec/ark/requests"
	"github.com/katasec/ark/resources/azure/cloudspaces"
	"gopkg.in/yaml.v2"
)

func PostCloudspace(s arkserver.Server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Deference pointer to get acsrepo and qclient
		db := s.GetAcsDb()
		qClient := s.GetCmdQ()

		// Set default acs name
		acsRequest := requests.CreateAzureCloudspaceRequest{
			Name: "default",
		}

		// Decode request body into acsRequest
		err := yaml.NewDecoder(r.Body).Decode(&acsRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get cloudspace from DB
		acs, err := db.GetCloudSpace(acsRequest.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error getting cloudspace: %s,", err)
			return
		}

		// Generate new cloudspace struct if none found in DB
		if acs.Name == "" {
			fmt.Println("Cloudspace not found in DB, creating new cloudspace")
			acs = *(cloudspaces.NewAzureCloudSpace())
			fmt.Println(acs.ToJson())
		}

		// Add environments from reuqest into struct
		for _, env := range acsRequest.Environments {
			acs.AddSpoke(env)
		}
		fmt.Println("The hub's name:", acs.Hub.Name)

		// Send request to queue
		subject := reflect.TypeOf(acsRequest).Name()
		log.Println("The subject is:" + subject)
		err = qClient.Send(subject, acs.ToJson()) // "azurecloudspace"
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Error: %s,", err)
			return
		} else {
			log.Println("Message sent successfully !")
		}

		// Save cloudspace to DB
		db.Create(&acs)
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
	})
}

func DeleteCloudspace(s arkserver.Server) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logx.LoggerFn().Info("In DeleteCloudspace handler")

		request := requests.DeleteAzureCloudspaceRequest{
			Name: "default",
		}

		err := yaml.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, request.ToYamlAzureCloudpace())

		qClient := s.GetCmdQ()

		qClient.Send("DeleteAzureCloudspaceRequest", request.ToJsonAzureCloudpace())

	})
}
