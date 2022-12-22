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
	arkRgName         = "rg-ark-001"
	arkStgAccountName = "arkstorage"
	arkSbNameSpace    = "ark"
)

// setupAzureDeps Sets up Azure Dependencies for Ark
func setupAzureDeps(ctx *pulumi.Context) error {

	// Create RG
	rg, err := resources.NewResourceGroup(ctx, arkRgName, &resources.ResourceGroupArgs{
		ResourceGroupName: pulumi.String(arkRgName),
	})
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

	// Add SB Namespace
	ns, err := servicebus.NewNamespace(ctx, arkSbNameSpace, &servicebus.NamespaceArgs{
		ResourceGroupName: rg.Name,
		Sku: servicebus.SBSkuArgs{
			Name: servicebus.SkuNameBasic,
			Tier: servicebus.SkuTierBasic,
		},
	})
	utils.ReturnError(err)
	ctx.Export("ns", ns.ServiceBusEndpoint)

	// Add Queue
	queue, err := servicebus.NewQueue(ctx, "command-queue", &servicebus.QueueArgs{
		ResourceGroupName:  rg.Name,
		EnablePartitioning: pulumi.Bool(true),
		NamespaceName:      ns.Name,
	})
	utils.ReturnError(err)

	ctx.Export("queueURN", queue.URN())
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
