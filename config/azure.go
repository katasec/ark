package config

type AzureConfig struct {
	ResourceGroupName string
	StorageConfig     AzureStorageConfig
	//MqConfig          AzureMqConfig
}

type AzureStorageConfig struct {
	StorageAccountName string
	//StorageEndpoint    string
	StorageKey string
	//LogsContainer        string
	PulumiStateContainer string
}

type AzureMqConfig struct {
	MqConnectionString string
	MqName             string
}
