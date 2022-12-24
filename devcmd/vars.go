package devcmd

import (
	"github.com/katasec/ark/utils"
	"github.com/katasec/ark/utils/docker"
)

var (
	// arkRgName         = "rg-ark-001"
	arkStgAccountName = "arkstorage"
	arkSbNameSpace    = "ark"

	// Used for checking prereqs
	checksPassed = true

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
)
