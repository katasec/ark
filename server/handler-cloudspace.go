package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	resources "github.com/katasec/ark/resources"
)

type CloudSpaceService interface {
	Create(c resources.CloudSpace) (id string, err error)
	Read(c resources.CloudSpace) (id string, err error)
	Update(c resources.CloudSpace) (id string, err error)
	Delete(c resources.CloudSpace) (err error)
}

type JsonCloudSpaceService struct {
}

func (s *JsonCloudSpaceService) Create(c resources.CloudSpace) (string, error) {
	return "", nil
}

func (s *Server) postCloudspace() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request := &resources.CloudSpace{}

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
