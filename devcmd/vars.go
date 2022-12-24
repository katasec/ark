package devcmd

import (
	"github.com/katasec/ark/utils"
	"github.com/katasec/ark/utils/docker"
)

var (

	// Used for checking prereqs
	checksPassed = true

	// Used for managing ark `dev` command
	d = NewDevCmd()

	// For PGSQL Docker
	DEV_PGSQL_IMAGE_NAME = "postgres:14.2-alpine"

	// Dev Instance defaults
	DevDbDefaultUser = "postgres"
	DevDbDefaultPass = "postgres"
	DevDbHost        = "127.0.0.1"
	DevDbPort        = 5432

	// Helps manage docker
	dh = docker.NewDockerHelper()

	// Spinner for status
	arkSpinner = utils.NewArkSpinner()

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
