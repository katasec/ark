package requests

import "github.com/katasec/ark/resources/v0/azure/cloudspaces"

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
