package dbtest

import (
	"encoding/json"

	"github.com/katasec/ark/sdk/v0/messages"
)

func genCloudSpace() string {
	referenceHubSubnets := []messages.SubnetsInfo{
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
	azureCloudspace := messages.AzureCloudspace{
		Name: "test",
		Hub: messages.VNETInfo{
			Name:          "hub",
			AddressPrefix: &hubPrefix,
			SubnetsInfo:   referenceHubSubnets,
		},
		Spokes: []messages.VNETInfo{
			{
				Name:          "spoke1",
				AddressPrefix: &spokePrefix,
				SubnetsInfo: []messages.SubnetsInfo{
					{
						Name:          "snet-test",
						Description:   "test net",
						AddressPrefix: "172.17.1.0/24",
					},
				},
			},
		},
	}

	jsonb, _ := json.MarshalIndent(azureCloudspace, "", "\t")
	return string(jsonb)
}
