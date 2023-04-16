package resources

type AzureCloudspace struct {
	Name     string
	Hub      VNETInfo
	Spokes   []VNETInfo
	Status   string
	Id       string
	UpdateId string
}

type VNETInfo struct {
	Name          string
	AddressPrefix *string
	SubnetsInfo   []SubnetsInfo
}

type SubnetsInfo struct {
	AddressPrefix string
	Description   string
	Name          string
	Tags          Tags
}

type Tags struct {
	Key   *string
	Value *string
}

var (
// Octet
)

func NewAzureCloudSpace(octet1 int, octet2 int) *AzureCloudspace {
	return &AzureCloudspace{}
}
