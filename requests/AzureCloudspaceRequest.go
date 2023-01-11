package requests

type AzureCloudspaceRequest struct {
	Env []Hub `json:"Env"`
	Hub Hub   `json:"Hub"`
}

type Hub struct {
	Name          string        `json:"Name"`
	AddressPrefix *string       `json:"AddressPrefix,omitempty"`
	SubnetsInfo   []SubnetsInfo `json:"SubnetsInfo,omitempty"`
}

type SubnetsInfo struct {
	AddressPrefix string `json:"AddressPrefix"`
	Description   string `json:"Description"`
	Name          string `json:"Name"`
	Tags          Tags   `json:"Tags"`
}

type Tags struct {
	Key   *string `json:"Key"`
	Value *string `json:"Value"`
}
