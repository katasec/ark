package messages

type AzureCloudspace struct {
	Action string     `json:"Action" yaml:"Action"`
	Name   string     `json:"Name" yaml:"Name"`
	Spokes []VNETInfo `json:"Spokes" yaml:"Spokes"`
	Hub    VNETInfo   `json:"Hub" yaml:"Hub"`
}

type VNETInfo struct {
	Name          string        `json:"Name" yaml:"Name"`
	AddressPrefix *string       `json:"AddressPrefix,omitempty" yaml:"AddressPrefix"`
	SubnetsInfo   []SubnetsInfo `json:"SubnetsInfo,omitempty" yaml:"SubnetsInfo"`
}

type SubnetsInfo struct {
	AddressPrefix string `json:"AddressPrefix" yaml:"AddressPrefix"`
	Description   string `json:"Description" yaml:"Description"`
	Name          string `json:"Name" yaml:"Name"`
	Tags          Tags   `json:"Tags" yaml:"Tags"`
}

type Tags struct {
	Key   *string `json:"Key" yaml:"Key"`
	Value *string `json:"Value" yaml:"Value"`
}
