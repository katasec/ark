package requests

import "github.com/katasec/ark/resources/azure/cloudspaces"

type DeleteAzureCloudspaceRequest struct {
	Name string `json:"Name" yaml:"Name"`
}

func (r *DeleteAzureCloudspaceRequest) ToJsonAzureCloudpace() string {
	acs := cloudspaces.NewAzureCloudSpace()
	acs.Name = r.Name
	return acs.ToJson()
}

func (r *DeleteAzureCloudspaceRequest) ToYamlAzureCloudpace() string {
	acs := cloudspaces.NewAzureCloudSpace()
	acs.Name = r.Name
	return acs.ToYaml()
}

func (r DeleteAzureCloudspaceRequest) GetResourceType() string {
	return "azurecloudspace"
}

func (r DeleteAzureCloudspaceRequest) GetActionType() string {
	return "delete"
}
