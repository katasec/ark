package cloudspaces

type SubnetsInfo struct {
	AddressPrefix string            `yaml:"AddressPrefix" `
	Description   string            `yaml:"Description"`
	Name          string            `yaml:"Name"`
	Tags          map[string]string `yaml:"Tags"`
}
