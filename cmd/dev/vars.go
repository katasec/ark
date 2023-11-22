package dev

import (
	dockerhelper "github.com/katasec/ark/docker-helper"
	"github.com/katasec/ark/utils"
)

var (

	// Used for checking prereqs
	checksPassed = true

	// Used for managing ark `dev` command
	d = NewDevCmd()

	// Docker Images
	DEV_PGSQL_IMAGE_NAME = "postgres:14.2-alpine"

	// DEV_ARK_SERVER_IMAGE_NAME   = "ghcr.io/katasec/arkserver:v0.0.7"
	// DEV_ARK_WORKER_IMAGE_NAME   = "ghcr.io/katasec/arkworker:v0.0.7"

	DEV_ARK_SERVER_IMAGE_NAME = "ghcr.io/katasec/arkserver:v0.0.11"
	DEV_ARK_WORKER_IMAGE_NAME = "ghcr.io/katasec/arkworker:v0.0.11"

	DEV_ARK_DB_MYSQL_IMAGE_NAME = "mysql:8.1.0"
	// Dev Instance defaults
	DevDbDefaultUser = "postgres"
	DevDbDefaultPass = "postgres"
	DevDbHost        = "127.0.0.1"
	DevDbPort        = 5432

	// Helps manage docker
	dh = dockerhelper.NewDockerHelper()

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
	StorageAccountName       = "StorageAccountName"
	StorageEndpoint          = "StorageEndpoint"
	StorageKey               = "StorageKey"
	PulumiStateContainerName = "PulumiStateContainerName"
	PulumiDbContainerName    = "PulumiDbContainerName"
)
