package cloudspaces

import "fmt"

func GenerateHubSubnets(octet1 int, octet2 int) []SubnetsInfo {
	prefix := fmt.Sprintf("%d.%d", octet1, octet2)

	return []SubnetsInfo{
		{
			Name:          "AzureFirewallSubnet",
			Description:   "Subnet for Azure Firewall",
			AddressPrefix: fmt.Sprintf("%s.0.0/26", prefix),
		},
		{
			Name:          "AzureBastionSubnet",
			Description:   "Subnet for Bastion Host",
			AddressPrefix: fmt.Sprintf("%s.0.64/26", prefix),
		},
		{
			Name:          "AzureFirewallManagementSubnet",
			Description:   "Subnet for Azure Firewall Mgmt",
			AddressPrefix: fmt.Sprintf("%s.0.128/26", prefix),
		},
		{
			Name:          "GatewaySubnet",
			Description:   "Subnet for VPN Gateway",
			AddressPrefix: fmt.Sprintf("%s.0.192/26", prefix),
		},
		{
			Name:          "snet-test",
			Description:   "Subnet for Testing purposes",
			AddressPrefix: fmt.Sprintf("%s.0.224/27", prefix),
		},
	}
}

func GenerateSpokeSubnets(octet1 int, octet2 int) []SubnetsInfo {
	prefix := fmt.Sprintf("%d.%d", octet1, octet2)

	return []SubnetsInfo{
		{
			Name:          "snet-tier1-agw",
			Description:   "Subnet for AGW",
			AddressPrefix: fmt.Sprintf("%s.1.0/24", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-agw",
			},
		},
		{
			Name:          "snet-tier1-webin",
			Description:   "Subnet for other LBs",
			AddressPrefix: fmt.Sprintf("%s.2.0/24", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-webin",
			},
		},
		{
			Name:          "snet-tier1-rsvd1",
			Description:   "Tier 1 reserved subnet",
			AddressPrefix: fmt.Sprintf("%s.3.0/25", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-rsvd1",
			},
		},
		{
			Name:          "snet-tier1-rsvd2",
			Description:   "Tier 1 reserved subnet",
			AddressPrefix: fmt.Sprintf("%s.3.128/25", prefix),
			Tags: map[string]string{
				"snet:role": "tier1-rsvd2",
			},
		},
		{
			Name:          "snet-tier2-wbapp",
			Description:   "Subnet for web apps",
			AddressPrefix: fmt.Sprintf("%s.4.0/23", prefix),
			Tags: map[string]string{
				"snet:role": "tier2-wbapp",
			},
		},
		{
			Name:          "snet-tier2-rsvd2",
			Description:   "Tier 2 reserved subnet",
			AddressPrefix: fmt.Sprintf("%s.6.0/24", prefix),
		},
		{
			Name:          "snet-tier2-pckr",
			Description:   "Subnet for packer",
			AddressPrefix: fmt.Sprintf("%s.7.0/24", prefix),
		},
		{
			Name:          "snet-tier2-vm",
			Description:   "Subnet for VMs",
			AddressPrefix: fmt.Sprintf("%s.8.0/21", prefix),
			Tags: map[string]string{
				"snet:role": "tier2-vm",
			},
		},
		{
			Name:          "snet-tier2-aks",
			Description:   "Subnet for AKS",
			AddressPrefix: fmt.Sprintf("%s.16.0/20", prefix),
		},
		{
			Name:          "snet-tier3-mi",
			Description:   "Subnet for managed instance",
			AddressPrefix: fmt.Sprintf("%s.32.0/26", prefix),
		},
		{
			Name:          "snet-tier3-dbaz",
			Description:   "Subnet for SQL Azure",
			AddressPrefix: fmt.Sprintf("%s.32.64/26", prefix),
		},
		{
			Name:          "snet-tier3-dblb",
			Description:   "Subnet for LB for SQL VM",
			AddressPrefix: fmt.Sprintf("%s.32.128/25", prefix),
		},
		{
			Name:          "snet-tier3-dbvm",
			Description:   "Subnet for SQL VM",
			AddressPrefix: fmt.Sprintf("%s.33.0/25", prefix),
		},
		{
			Name:          "snet-tier3-strg",
			Description:   "Subnet for storage account/fileshares",
			AddressPrefix: fmt.Sprintf("%s.33.128/25", prefix),
		},
		{
			Name:          "snet-tier3-redis",
			Description:   "Subnet for redis cache",
			AddressPrefix: fmt.Sprintf("%s.34.0/25", prefix),
		},
	}
}
