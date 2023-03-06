package messages

type AzureManagedCluster struct {
	VnetResourceGroup string
	SubNetName        string
	VnetName          string
	Aks               Aks
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
