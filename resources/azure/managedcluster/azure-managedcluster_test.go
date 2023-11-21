package resources

import (
	"testing"

	yaml "github.com/katasec/utils/yaml"
)

func TestManagedCluster(t *testing.T) {
	cluster := NewAzureManagedCluster("rg-ark", "vnet-ark", "subnet-ark", "aks-ark")
	x, err := yaml.YamlMarshall[*AzureManagedCluster](cluster)
	if err != nil {
		t.Error(err)
	}
	t.Log("\n" + x)
}
