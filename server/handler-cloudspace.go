package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	resources "github.com/katasec/ark/resources"
)

type CloudSpaceService interface {
	Create(c resources.CloudspaceRequest) (id string, err error)
	Delete(c resources.CloudspaceRequest) (err error)
}

type JsonCloudSpaceService struct {
}

func (s *JsonCloudSpaceService) Create(c resources.CloudspaceRequest) (string, error) {
	return "", nil
}

type JsonDB struct {
	cloudspaces []resources.CloudspaceRequest
}

func JsonRepository(c resources.CloudspaceRequest) {

}

func (s *Server) postCloudspace() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request := &resources.CloudspaceRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Cloud Space: %+v", request)

	})
}

func DecodingError(w http.ResponseWriter, err error, reqName string) {
	errMsg := fmt.Sprintf("HTTP Error: %s. Error decoding request as a %s request", err.Error(), reqName)
	http.Error(w, errMsg, http.StatusBadRequest)

}
