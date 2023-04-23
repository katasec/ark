package cloudspaces

type VNETInfo struct {
	Name          string        `yaml:"Name"`
	AddressPrefix string        `yaml:"AddressPrefix"`
	SubnetsInfo   []SubnetsInfo `yaml:"SubnetsInfo"`
}
