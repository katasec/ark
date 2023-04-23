package requests

import "github.com/katasec/ark/resources/v0/azure/cloudspaces"

type AzureCloudspaceRequest struct {
	Kind         string
	Environments []string
}

func (r *AzureCloudspaceRequest) ToAzureCloudpace() string {
	acs := cloudspaces.NewAzureCloudSpace()
	for _, env := range r.Environments {
		acs.AddSpoke(env)
	}
	return acs.ToJson()
}

// type Hub struct {
// 	Name          string        `json:"Name"`
// 	AddressPrefix *string       `json:"AddressPrefix,omitempty"`
// 	SubnetsInfo   []SubnetsInfo `json:"SubnetsInfo,omitempty"`
// }

// type SubnetsInfo struct {
// 	AddressPrefix string `json:"AddressPrefix"`
// 	Description   string `json:"Description"`
// 	Name          string `json:"Name"`
// 	Tags          Tags   `json:"Tags"`
// }

// type Tags struct {
// 	Key   *string `json:"Key"`
// 	Value *string `json:"Value"`
// }
