package dbtest

import (
	"encoding/json"

	"github.com/katasec/ark/resources/v0/azure/cloudspaces"
)

func genCloudSpace() cloudspaces.AzureCloudspace {
	referenceHubSubnets := []cloudspaces.SubnetsInfo{
		{
			Name:          "AzureFirewallSubnet",
			Description:   "Subnet for Azure Firewall",
			AddressPrefix: "172.16.0.0/26",
		},
		{
			Name:          "AzureBastionSubnet",
			Description:   "Subnet for Bastion",
			AddressPrefix: "172.16.0.64/26",
		},
		{
			Name:          "AzureFirewallManagementSubnet",
			Description:   "Subnet for VPN Gateway",
			AddressPrefix: "172.16.0.128/26",
		},
		{
			Name:          "GatewaySubnet",
			Description:   "Subnet for VPN Gateway",
			AddressPrefix: "172.16.0.192/27",
		},
		{
			Name:          "snet-test",
			Description:   "Subnet for VPN Gateway",
			AddressPrefix: "172.16.0.224/27",
		},
	}

	hubPrefix := "172.16.0.0/24"
	spokePrefix := "172.17.0.0/16"
	azureCloudspace := cloudspaces.AzureCloudspace{
		Name: "test",
		Hub: cloudspaces.VNETInfo{
			Name:          "hub",
			AddressPrefix: hubPrefix,
			SubnetsInfo:   referenceHubSubnets,
		},
		Spokes: []cloudspaces.VNETInfo{
			{
				Name:          "spoke1",
				AddressPrefix: spokePrefix,
				SubnetsInfo: []cloudspaces.SubnetsInfo{
					{
						Name:          "snet-test",
						Description:   "test net",
						AddressPrefix: "172.17.1.0/24",
					},
				},
			},
		},
	}

	return azureCloudspace
}

func genCloudSpaceJson() string {
	acs := genCloudSpace()

	jsonb, _ := json.MarshalIndent(acs, "", "\t")
	return string(jsonb)
}
