package server

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	logx "github.com/katasec/ark/log"
	"github.com/katasec/tableio"

	"github.com/katasec/ark/requests"
	"github.com/katasec/ark/resources/azure/cloudspaces"
	"gopkg.in/yaml.v2"
)

func (s *Server) PostCloudspace() http.HandlerFunc {

	acsTable, err := tableio.NewTableIO[cloudspaces.AzureCloudspace](s.config.DbConfig.DriverName, s.config.DbConfig.DataSourceName)
	if err != nil {
		log.Printf("Error creating tableio struct: %s\n", err)
		return nil
	} else {
		log.Printf("Created tableio struct: ")
	}
	acsTable.CreateTableIfNotExists()
	if err != nil {
		log.Printf("Error creating table %s\n", err)
		return nil
	} else {
		log.Printf("Created table\n")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
		rows, err := acsTable.ByName(acsRequest.Name)
		log.Println("looking for cloudspace:" + acsRequest.Name)
		log.Println("Number of rows returned:" + string(len(rows)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error getting cloudspace: %s,", err)
			return
		}

		var acs cloudspaces.AzureCloudspace
		// Generate new cloudspace struct if none found in DB
		if len(rows) == 0 {
			fmt.Println("Cloudspace not found in DB, creating new cloudspace")
			acs = *(cloudspaces.NewAzureCloudSpace())
			fmt.Println(acs.ToJson())
		} else {
			acs = rows[0]
		}

		// Add environments from reuqest into struct
		for _, env := range acsRequest.Environments {
			acs.AddSpoke(env)
		}
		fmt.Println("The hub's name:", acs.Hub.Name)

		// Send request to queue
		subject := reflect.TypeOf(acsRequest).Name()
		log.Println("The subject is:" + subject)
		err = s.cmdQ.Send(subject, acs.ToJson()) // "azurecloudspace"
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Error: %s,", err)
			return
		} else {
			log.Println("Message sent successfully !")
		}

		// Save cloudspace to DB
		acsTable.Insert(acs)

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
	})

}

func (s *Server) DeleteCloudspace() http.HandlerFunc {

	fmt.Println("In DeleteCloudspace() http.HandlerFunc")
	acsTable, err := tableio.NewTableIO[cloudspaces.AzureCloudspace](s.config.DbConfig.DriverName, s.config.DbConfig.DataSourceName)
	if err != nil {
		log.Printf("Error creating tableio struct: %s\n", err)
		return nil
	}
	acsTable.CreateTableIfNotExists()
	if err != nil {
		log.Printf("Error creating table: %s\n", err)
		return nil
	}

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

		err = s.cmdQ.Send("DeleteAzureCloudspaceRequest", request.ToJsonAzureCloudpace())
		if err != nil {
			log.Printf("Error sending message to queue: %s\n", err.Error())
		}

	})
}
