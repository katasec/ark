package dev

import (
	"fmt"
	"os"
	"strings"

	"github.com/katasec/ark/utils"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/servicebus"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	arkResourceGroup  *resources.ResourceGroup
	arkRgName         = "rg-ark-001"
	arkStgAccountName = "arkstorage"
	arkSbNameSpace    = "ark"
)

func createRgFunc(ctx *pulumi.Context) error {
	rg, err := resources.NewResourceGroup(ctx, arkRgName, &resources.ResourceGroupArgs{
		ResourceGroupName: pulumi.String(arkRgName),
	})
	utils.ReturnError(err)

	ctx.Export("rgName", rg.Name)

	rg.Name.ApplyT(func(name string) string {
		f, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE, 0600)
		fmt.Fprintf(f, "RG Name:"+name)
		return name
	})

	return nil
}

func createStrgFunc(ctx *pulumi.Context) error {

	account, err := storage.NewStorageAccount(ctx, arkStgAccountName, &storage.StorageAccountArgs{
		ResourceGroupName: pulumi.String(arkRgName),
		AccessTier:        storage.AccessTierHot,
		Sku: &storage.SkuArgs{
			Name: pulumi.String(storage.SkuName_Standard_LRS),
		},
		Kind: pulumi.String(storage.KindStorageV2),
	})
	utils.ReturnError(err)

	ctx.Export("accountName", account.Name)

	account.Name.ApplyT(func(name string) string {
		f, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE, 0600)
		fmt.Fprintf(f, "Account Name:"+name)
		return name
	})
	return nil
}

func createSbNsFunc(ctx *pulumi.Context) error {
	ns, err := servicebus.NewNamespace(ctx, arkSbNameSpace, &servicebus.NamespaceArgs{
		ResourceGroupName: pulumi.String(arkRgName),
		Sku: servicebus.SBSkuArgs{
			Name: servicebus.SkuNameBasic,
			Tier: servicebus.SkuTierBasic,
		},
	})
	utils.ReturnError(err)

	ctx.Export("all", ns)
	ctx.Export("ns", ns.Name)

	ns.Name.ApplyT(func(name string) string {
		f, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE, 0600)
		fmt.Fprintf(f, "Sb Name:"+name)
		return name
	})
	return nil
}

func createSbQueue(ctx *pulumi.Context) error {

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
