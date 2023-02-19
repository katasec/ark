package messages

type AzureAks struct {
	VnetResourceGroup string
	SubNetName        string
	VnetName          string

	Aks AksConfig
}

type AksConfig struct {
	Name                 string
	ResourceGroup        string
	ServicePrincipal     string
	VmSize               string
	DnsPrefix            string
	EnablePrivateCluster bool
	NetworkProfile       NetworkProfile
}

type NetworkProfile struct {
	NetworkPlugin    string
	NetworkPolicy    string
	DockerBridgeCidr string
}
