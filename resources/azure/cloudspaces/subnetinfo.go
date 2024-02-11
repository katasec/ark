package cloudspaces

type SubnetsInfo struct {
	Name               string `yaml:"Name"`
	ResourceGroupName  string `yaml:"ResourceGroupName" `
	VirtualNetworkName string `yaml:"VirtualNetworkName" `
	AddressPrefixes    string `yaml:"AddressPrefixes" `

	Description string            `yaml:"Description"`
	Tags        map[string]string `yaml:"Tags"`
}
