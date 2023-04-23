package requests

import "github.com/katasec/ark/resources/v0/azure/cloudspaces"

type CreateAzureCloudspaceRequest struct {
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
