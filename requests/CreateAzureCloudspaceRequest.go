package requests

import (
	"github.com/katasec/ark/resources/azure/cloudspaces"
)

type CreateAzureCloudspaceRequest struct {
	Name         string
	Kind         string
	Environments []string
}

func (r *CreateAzureCloudspaceRequest) ToJsonAzureCloudpace() string {
	acs := cloudspaces.NewAzureCloudSpace()
	for _, env := range r.Environments {
		acs.AddSpoke(env)
	}
	return acs.ToJson()
}

func (r *CreateAzureCloudspaceRequest) ToYamlAzureCloudpace() string {
	acs := cloudspaces.NewAzureCloudSpace()
	for _, env := range r.Environments {
		acs.AddSpoke(env)
	}
	return acs.ToYaml()
}

// GetResourceType For e.g. "azurecloudspace"
func (r CreateAzureCloudspaceRequest) GetResourceType() string {
	return "azurecloudspace"
}

// GetRequestType For e.g. "create"
func (r CreateAzureCloudspaceRequest) GetActionType() string {
	return "create"
}
