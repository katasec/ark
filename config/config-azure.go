package config

type AzureConfig struct {
	ResourceGroupName string
	StorageConfig     AzureLogConfig
	MqConfig          AzureMqConfig
}

type AzureLogConfig struct {
	LogStorageAccountName string
	LogStorageEndpoint    string
	LogStorageKey         string
	LogsContainer         string
	PulumiStateContainer  string
}

type AzureMqConfig struct {
	MqConnectionString string
	MqName             string
}
