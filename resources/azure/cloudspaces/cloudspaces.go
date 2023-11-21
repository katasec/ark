package cloudspaces

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type AzureCloudspace struct {
	Name     string     `yaml:"Name"`
	Hub      VNETInfo   `yaml:"Hub"`
	Spokes   []VNETInfo `yaml:"Spokes"`
	Status   string     `yaml:"Status"`
	Id       string
	UpdateId string

	hubOctet1 int
	hubOctet2 int

	firstOctet2 int
}

// NewAzureCloudSpace Create a new Azure Cloudspace
func NewAzureCloudSpace(oct1 ...int) *AzureCloudspace {

	// variables DefaultOctet1 and DefaultOctet2  default value for the first/second octets
	octet1 := DefaultOctet1
	if len(oct1) > 0 {
		octet1 = oct1[0]
	}

	return &AzureCloudspace{
		Name: "default",
		Hub: VNETInfo{
			Name:          "vnet-hub",
			AddressPrefix: fmt.Sprintf("%d.%d.0.0/16", octet1, DefaultOctet2),
			SubnetsInfo:   GenerateHubSubnets(octet1, DefaultOctet2),
		},
		hubOctet2:   DefaultOctet2,
		firstOctet2: DefaultOctet2 + 1,
		hubOctet1:   octet1,
	}
}

// AddSpoke Add a spoke to the cloudspace
func (acs *AzureCloudspace) AddSpoke(name string) error {

	// Return if the spoke already exists
	if acs.IsSpoke(name) {
		return nil
	} else if len(acs.Spokes) >= len(AllOctet2) {
		// Return if the number of spokes is greater than allowed
		return fmt.Errorf("can't have more than %d spokes", len(AllOctet2))
	}

	// Determine 2nd octet for the new spoke
	octet2 := acs.NextAvailableOctet2()

	// Create a new spoke
	newSpoke := VNETInfo{
		Name:          fmt.Sprintf("%s%s", VnetPrefix, name),
		AddressPrefix: fmt.Sprintf("%d.%d.0.0/16", acs.hubOctet1, octet2),
		SubnetsInfo:   GenerateSpokeSubnets(acs.hubOctet1, octet2),
	}

	// Add the new spoke to the list of spokes
	acs.Spokes = append(acs.Spokes, newSpoke)

	return nil
}

// UsedOctet2s Get a list of used 2nd Octets
func (acs *AzureCloudspace) UsedOctet2s() []int {
	// Get List of used octets
	usedOctets := []int{}
	for _, spoke := range acs.Spokes {
		octets := strings.Split(spoke.AddressPrefix, ".")
		octet, _ := strconv.Atoi(octets[1]) // 2nd octet
		usedOctets = append(usedOctets, octet)
	}

	return usedOctets
}

// NextAvailableOctet2 Get the next available 2nd Octet
func (acs *AzureCloudspace) NextAvailableOctet2() int {
	used := acs.UsedOctet2s()
	for _, octet := range AllOctet2 {
		if !contains(used, octet) {
			return octet
		}
	}

	return 0
}

// contains returns true if the slice contains the element
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SpokeCount Get the number of spokes
func (acs *AzureCloudspace) SpokeCount() int {
	return len(acs.Spokes)
}

// IsSpoke Check if the spoke exists
func (acs *AzureCloudspace) IsSpoke(name string) bool {

	vnetName := fmt.Sprintf("%s%s", VnetPrefix, name)
	for _, spoke := range acs.Spokes {
		if strings.EqualFold(spoke.Name, vnetName) {
			return true
		}
	}

	return false
}

func (acs *AzureCloudspace) ToJson() string {
	b, err := json.Marshal(acs)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(b)
}

func (acs *AzureCloudspace) ToYaml() string {
	b, err := yaml.Marshal(acs)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(b)
}
