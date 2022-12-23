package devcmd

var (

	// Pulumi Stack Details
	ProjectNamePrefix = "ark-init"
	StackName         = "dev"

	// Azure Resources Names
	ResourceGroupPrefix = "rg-ark-"

	// Pulumi Export Names
	ResourceGroupName        = "ResourceGroupName"
	StgAccountPrefix         = "arkstorage"
	AsbNsPrefix              = "arkns"
	PrimaryStorageKey        = "PrimaryStorageKey"
	MqConnectionString       = "MqConnectionString"
	CommandQueueName         = "CommandQueueName"
	LogStorageAccountName    = "LogStorageAccountName"
	LogStorageEndpoint       = "LogStorageEndpoint"
	LogStorageKey            = "LogStorageKey"
	LogContainerName         = "LogContainerName"
	PulumiStateContainerName = "PulumiStateContainerName"
)
