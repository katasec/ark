package resources

type AzureManagedCluster struct {
	VnetResourceGroup string
	SubNetName        string
	VnetName          string
	Aks               Aks
}

func NewAzureManagedCluster(vnetResourceGroup string, vnetName string, subnetName string, aksName string) *AzureManagedCluster {
	return &AzureManagedCluster{
		VnetResourceGroup: vnetResourceGroup,
		VnetName:          vnetName,
		SubNetName:        subnetName,
		Aks: Aks{
			Name:             aksName,
			ResourceGroup:    "rg-ark-" + aksName,
			ServicePrincipal: "sp-ark-" + aksName,
			NetworkProfile: NetworkProfile{
				NetworkPlugin:    "azure",
				NetworkPolicy:    "calico",
				DockerBridgeCIDR: "172.17.0.0/16",
			},
		},
	}
}

type Aks struct {
	Name                 string
	ResourceGroup        string
	ServicePrincipal     string
	VMSize               string
	EnablePrivateCluster bool
	NetworkProfile       NetworkProfile
}

type NetworkProfile struct {
	NetworkPlugin    string
	NetworkPolicy    string
	DockerBridgeCIDR string
}
