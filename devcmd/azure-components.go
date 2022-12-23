package devcmd

import (
	"fmt"
	"strings"

	"github.com/katasec/ark/utils"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/servicebus"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	// arkRgName         = "rg-ark-001"
	arkStgAccountName = "arkstorage"
	arkSbNameSpace    = "ark"
)

// setupAzureComponents Sets up Azure Components for Ark
func setupAzureComponents(ctx *pulumi.Context) error {

	// Create RG
	rg, err := resources.NewResourceGroup(ctx, ResourceGroupPrefix, nil)
	utils.ReturnError(err)
	ctx.Export("rgName", rg.Name)

	// Add Storage Account
	account, err := storage.NewStorageAccount(ctx, arkStgAccountName, &storage.StorageAccountArgs{
		ResourceGroupName: rg.Name,
		AccessTier:        storage.AccessTierHot,
		Sku: &storage.SkuArgs{
			Name: pulumi.String(storage.SkuName_Standard_LRS),
		},
		Kind: pulumi.String(storage.KindStorageV2),
	})
	utils.ReturnError(err)
	ctx.Export("accountName", account.Name)

	// Get Storage Key
	ctx.Export("primaryStorageKey", pulumi.All(rg.Name, account.Name).ApplyT(func(args []interface{}) (string, error) {
		accountKeys, err := storage.ListStorageAccountKeys(ctx, &storage.ListStorageAccountKeysArgs{
			ResourceGroupName: args[0].(string),
			AccountName:       args[1].(string),
		})
		if err != nil {
			return "", err
		}

		return accountKeys.Keys[0].Value, nil
	}))

	// Create ASB Namespace
	ns, err := servicebus.NewNamespace(ctx, arkSbNameSpace, &servicebus.NamespaceArgs{
		ResourceGroupName: rg.Name,
		Sku: servicebus.SBSkuArgs{
			Name: servicebus.SkuNameBasic,
			Tier: servicebus.SkuTierBasic,
		},
	})
	utils.ReturnError(err)
	ctx.Export("ns", ns.ServiceBusEndpoint)

	// Create New Auth Rule
	authRuleName := "ReadWrite"
	authRule, err := servicebus.NewNamespaceAuthorizationRule(ctx, "ReadWrite", &servicebus.NamespaceAuthorizationRuleArgs{
		AuthorizationRuleName: pulumi.String(authRuleName),
		NamespaceName:         ns.Name,
		ResourceGroupName:     rg.Name,
		Rights: servicebus.AccessRightsArray{
			servicebus.AccessRightsListen,
			servicebus.AccessRightsSend,
		},
	})
	utils.ReturnError(err)

	// Export Connection String
	ctx.Export("ASBPrimaryConnectionString", pulumi.All(rg.Name, ns.Name, authRule.Name).ApplyT(func(args []interface{}) (string, error) {
		keys, err := servicebus.ListNamespaceKeys(ctx, &servicebus.ListNamespaceKeysArgs{
			ResourceGroupName:     args[0].(string),
			NamespaceName:         args[1].(string),
			AuthorizationRuleName: authRuleName,
		})
		if err != nil {
			return "", err
		}

		return keys.PrimaryConnectionString, nil
	}))

	// Create Queue in ASB namespace
	queue, err := servicebus.NewQueue(ctx, "command-queue", &servicebus.QueueArgs{
		ResourceGroupName:  rg.Name,
		EnablePartitioning: pulumi.Bool(true),
		NamespaceName:      ns.Name,
	})
	utils.ReturnError(err)
	ctx.Export("queueName", queue.Name)

	fmt.Println("End")
	return nil
}

func getReference(stackFQDN string, key string) (output string, err error) {
	myCmd := fmt.Sprintf("pulumi stack -s %s output %s", stackFQDN, key)

	value, err := utils.ExecShellCmd(myCmd)
	value = strings.Trim(value, "\n")

	return value, err
}

func getDefaultPulumiOrg() (string, error) {

	value, err := utils.ExecShellCmd("pulumi org get-default")
	value = strings.Trim(value, "\n")

	return value, err
}
