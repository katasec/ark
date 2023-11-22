package dev

import (
	"github.com/katasec/ark/utils"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// setupAzureComponents Sets up Azure Components for Ark
func setupAzureComponents(ctx *pulumi.Context) error {

	// Create RG
	rg, err := resources.NewResourceGroup(ctx, ResourceGroupPrefix, &resources.ResourceGroupArgs{})
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
	ctx.Export(StorageEndpoint, account.PrimaryEndpoints.Blob())
	ctx.Export(StorageAccountName, account.Name)

	// Get Storage Key
	ctx.Export(StorageKey, pulumi.All(rg.Name, account.Name).ApplyT(func(args []interface{}) (string, error) {
		accountKeys, err := storage.ListStorageAccountKeys(ctx, &storage.ListStorageAccountKeysArgs{
			ResourceGroupName: args[0].(string),
			AccountName:       args[1].(string),
		})
		if err != nil {
			return "", err
		}

		return accountKeys.Keys[0].Value, nil
	}))

	// Pulumi state blob container
	pulumiStateContainer, err := storage.NewBlobContainer(ctx, "pulumistate", &storage.BlobContainerArgs{
		AccountName:       account.Name,
		ResourceGroupName: rg.Name,
	})
	utils.ReturnError(err)
	ctx.Export(PulumiStateContainerName, pulumiStateContainer.Name)

	return nil
}
