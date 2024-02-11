package cloudspaces

type VNETInfo struct {
	ResourceGroup string
	Name          string        `yaml:"Name" json:"Name"`
	AddressPrefix string        `yaml:"AddressPrefix" json:"AddressPrefix"`
	SubnetsInfo   []SubnetsInfo `yaml:"SubnetsInfo" json:"SubnetsInfo"`
}
