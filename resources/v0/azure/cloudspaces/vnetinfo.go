package cloudspaces

type VNETInfo struct {
	Name          string
	AddressPrefix string
	SubnetsInfo   []SubnetsInfo
}
