package cloudspaces

import "fmt"

func GenerateHubSubnets(octet1 int, octet2 int, vnetName string, resourceGroupName string) []SubnetsInfo {
	prefix := fmt.Sprintf("%d.%d", octet1, octet2)

	return []SubnetsInfo{
		{
			VirtualNetworkName: vnetName,
			Name:               "AzureFirewallSubnet",
			Description:        "Subnet for Azure Firewall",
			AddressPrefixes:    fmt.Sprintf("%s.0.0/26", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "AzureBastionSubnet",
			Description:        "Subnet for Bastion Host",
			AddressPrefixes:    fmt.Sprintf("%s.0.64/26", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "AzureFirewallManagementSubnet",
			Description:        "Subnet for Azure Firewall Mgmt",
			AddressPrefixes:    fmt.Sprintf("%s.0.128/26", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "GatewaySubnet",
			Description:        "Subnet for VPN Gateway",
			AddressPrefixes:    fmt.Sprintf("%s.0.192/26", prefix),
			ResourceGroupName:  resourceGroupName,
		},
	}
}

func GenerateSpokeSubnets(octet1 int, octet2 int, vnetName string, resourceGroupName string) []SubnetsInfo {
	prefix := fmt.Sprintf("%d.%d", octet1, octet2)

	return []SubnetsInfo{
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier1-agw",
			Description:        "Subnet for AGW",
			AddressPrefixes:    fmt.Sprintf("%s.1.0/24", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-agw",
			},
			ResourceGroupName: resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier1-webin",
			Description:        "Subnet for other LBs",
			AddressPrefixes:    fmt.Sprintf("%s.2.0/24", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-webin",
			},
			ResourceGroupName: resourceGroupName,
		},

		{
			Name:            "snet-tier1-rsvd1",
			Description:     "Tier 1 reserved subnet",
			AddressPrefixes: fmt.Sprintf("%s.3.0/25", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-rsvd1",
			},
			ResourceGroupName: resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier1-rsvd2",
			Description:        "Tier 1 reserved subnet",
			AddressPrefixes:    fmt.Sprintf("%s.3.128/25", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-rsvd2",
			},
			ResourceGroupName: resourceGroupName,
		},
		{
			Name:            "snet-tier2-wbapp",
			Description:     "Subnet for web apps",
			AddressPrefixes: fmt.Sprintf("%s.4.0/23", prefix),
			Tags: map[string]string{
				"snet:role": "tier2-wbapp",
			},
			ResourceGroupName: resourceGroupName,
		},
		{
			Name:              "snet-tier2-rsvd2",
			Description:       "Tier 2 reserved subnet",
			AddressPrefixes:   fmt.Sprintf("%s.6.0/24", prefix),
			ResourceGroupName: resourceGroupName,
		},
		{
			Name:              "snet-tier2-pckr",
			Description:       "Subnet for packer",
			AddressPrefixes:   fmt.Sprintf("%s.7.0/24", prefix),
			ResourceGroupName: resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier2-vm",
			Description:        "Subnet for VMs",
			AddressPrefixes:    fmt.Sprintf("%s.8.0/21", prefix),
			Tags: map[string]string{
				"snet:role": "tier2-vm",
			},
			ResourceGroupName: resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier2-aks",
			Description:        "Subnet for AKS",
			AddressPrefixes:    fmt.Sprintf("%s.16.0/20", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier3-mi",
			Description:        "Subnet for managed instance",
			AddressPrefixes:    fmt.Sprintf("%s.32.0/26", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier3-dbaz",
			Description:        "Subnet for SQL Azure",
			AddressPrefixes:    fmt.Sprintf("%s.32.64/26", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier3-dblb",
			Description:        "Subnet for LB for SQL VM",
			AddressPrefixes:    fmt.Sprintf("%s.32.128/25", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier3-dbvm",
			Description:        "Subnet for SQL VM",
			AddressPrefixes:    fmt.Sprintf("%s.33.0/25", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier3-strg",
			Description:        "Subnet for storage account/fileshares",
			AddressPrefixes:    fmt.Sprintf("%s.33.128/25", prefix),
			ResourceGroupName:  resourceGroupName,
		},
		{
			VirtualNetworkName: vnetName,
			Name:               "snet-tier3-redis",
			Description:        "Subnet for redis cache",
			AddressPrefixes:    fmt.Sprintf("%s.34.0/25", prefix),
			ResourceGroupName:  resourceGroupName,
		},
	}
}
