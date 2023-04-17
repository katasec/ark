package cloudspaces

import "fmt"

type AzureCloudspace struct {
	Name     string
	Hub      VNETInfo
	Spokes   []VNETInfo
	Status   string
	Id       string
	UpdateId string
}

func NewAzureCloudSpace(octet1 ...int) *AzureCloudspace {

	// DefaultOctet1 is the default value for the first octet of the address prefix
	myOctet1 := DefaultOctet1
	if len(octet1) > 0 {
		myOctet1 = octet1[0]
	}

	// DefaultOctet2 is the default value for the second octet of the address prefix
	return &AzureCloudspace{
		Hub: VNETInfo{
			Name:          "vnet-hub",
			AddressPrefix: fmt.Sprintf("%d.%d.0.0/24", myOctet1, DefaultOctet2),
		},
	}
}

func (acs *AzureCloudspace) AddSpoke(name string) {

	acs.Spokes = append(acs.Spokes, VNETInfo{
		Name:          fmt.Sprintf("vnet-%s", name),
		AddressPrefix: fmt.Sprintf("%d.%d.0.0/24", DefaultOctet1, Octet2Range[len(acs.Spokes)]),
	})
}

func (acs *AzureCloudspace) IsSpoke(name string) bool {
	return false
}
