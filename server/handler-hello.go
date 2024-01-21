package server

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/katasec/ark/requests"
	"github.com/katasec/ark/resources/azure/hello"
	"github.com/katasec/tableio"
	"gopkg.in/yaml.v2"
)

func (s *Server) PostHello() http.HandlerFunc {

	log.Println("PostHello()")
	resourceTable, err := tableio.NewTableIO[hello.Hello](s.config.DbConfig.DriverName, s.config.DbConfig.DataSourceName)
	if err != nil {
		log.Printf("Error creating tableio struct: %s\n", err)
		return nil
	} else {
		log.Printf("Created tableio struct: ")
	}

	resourceTable.CreateTableIfNotExists()
	if err != nil {
		log.Printf("Error creating table %s\n", err)
		return nil
	} else {
		log.Printf("Created table\n")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set default acs name
		request := requests.CreateHelloRequest{}

		// Decode request body into HelloRequest
		err := yaml.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get cloudspace from DB
		rows, err := resourceTable.ByName(request.Name)
		log.Println("looking for  hello:" + request.Name)
		log.Printf("Number of rows returned:%d\n", len(rows))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error getting cloudspace: %s,", err)
			return
		}

		var resource hello.Hello
		// Generate new cloudspace struct if none found in DB
		if len(rows) == 0 {
			fmt.Println("Cloudspace not found in DB, creating new cloudspace")
			resource = hello.Hello{}
			fmt.Println(resource.ToJson())
		} else {
			resource = rows[0]
		}

		// Send request to queue
		subject := reflect.TypeOf(request).Name()
		err = s.cmdQ.Send(subject, resource.ToJson())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Error: %s,", err)
			return
		} else {
			log.Println("Message sent successfully !")
		}

		// Save cloudspace to DB
		resourceTable.Insert(resource)

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
	})

}
