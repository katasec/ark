package requests

import "errors"

type RequestInterface interface {
	// GetResourceType For e.g. "azurecloudspace
	GetResourceType() string

	// GetRequestType For e.g. "create"
	GetRequestType() string
}
type Request any
type Requests struct{}

var (
	requests = map[string]Request{
		"createazurecloudspace": CreateAzureCloudspaceRequest{},
		"deleteazurecloudspace": DeleteAzureCloudspaceRequest{},
	}
)

func (r Requests) Create(typeName string) (interface{}, error) {
	if _, ok := requests[typeName]; ok {
		return requests[typeName], nil
	} else {
		return nil, errors.New("Request type not found")
	}
}
