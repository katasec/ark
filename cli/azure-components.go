package cli

import (
	"github.com/katasec/ark/utils"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/servicebus"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// setupAzureComponents Sets up Azure Components for Ark
func setupAzureComponents(ctx *pulumi.Context) error {

	// Create RG
	rg, err := resources.NewResourceGroup(ctx, ResourceGroupPrefix, nil)
	utils.ReturnError(err)
	ctx.Export(ResourceGroupName, rg.Name)

	// Add Storage Account
	account, err := storage.NewStorageAccount(ctx, StgAccountPrefix, &storage.StorageAccountArgs{
		ResourceGroupName: rg.Name,
		AccessTier:        storage.AccessTierHot,
		Sku: &storage.SkuArgs{
			Name: pulumi.String(storage.SkuName_Standard_LRS),
		},
		Kind: pulumi.String(storage.KindStorageV2),
	})
	utils.ReturnError(err)
	ctx.Export(LogStorageEndpoint, account.PrimaryEndpoints.Blob())
	ctx.Export(LogStorageAccountName, account.Name)

	// Get Storage Key
	ctx.Export(LogStorageKey, pulumi.All(rg.Name, account.Name).ApplyT(func(args []interface{}) (string, error) {
		accountKeys, err := storage.ListStorageAccountKeys(ctx, &storage.ListStorageAccountKeysArgs{
			ResourceGroupName: args[0].(string),
			AccountName:       args[1].(string),
		})
		if err != nil {
			return "", err
		}

		return accountKeys.Keys[0].Value, nil
	}))

	// Storage Containers
	arklogsContainer, err := storage.NewBlobContainer(ctx, "arklogs", &storage.BlobContainerArgs{
		AccountName:       account.Name,
		ResourceGroupName: rg.Name,
	})
	utils.ReturnError(err)
	ctx.Export(LogContainerName, arklogsContainer.Name)

	// Storage Containers
	pulumiStateContainer, err := storage.NewBlobContainer(ctx, "pulumistate", &storage.BlobContainerArgs{
		AccountName:       account.Name,
		ResourceGroupName: rg.Name,
	})
	utils.ReturnError(err)
	ctx.Export(PulumiStateContainerName, pulumiStateContainer.Name)

	// Create ASB Namespace
	ns, err := servicebus.NewNamespace(ctx, AsbNsPrefix, &servicebus.NamespaceArgs{
		ResourceGroupName: rg.Name,
		Sku: servicebus.SBSkuArgs{
			Name: servicebus.SkuNameBasic,
			Tier: servicebus.SkuTierBasic,
		},
	})
	utils.ReturnError(err)

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
	ctx.Export(MqConnectionString, pulumi.All(rg.Name, ns.Name, authRule.Name).ApplyT(func(args []interface{}) (string, error) {
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
	ctx.Export(CommandQueueName, queue.Name)

	return nil
}
