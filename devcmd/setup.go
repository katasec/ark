package devcmd

import (
	"github.com/katasec/ark/utils/docker"
)

var (
	d = NewDevCmd()

	dh = docker.NewDockerHelper()

	DEV_PGSQL_IMAGE_NAME = "postgres:14.2-alpine"

	// Dev Instance defaults
	DevDbDefaultUser = "postgres"
	DevDbDefaultPass = "postgres"
	DevDbHost        = "127.0.0.1"
	DevDbPort        = 5432
)

func Setup() {
	// Create Cloud resources with Pulumi
	d.Setup()

	setupDb()
}

func SetupDelete() {
	// Delete Cloud resources with Pulumi
	d.Delete()
}

func setupDb() {

	// Start Postgres Container
	imageName := DEV_PGSQL_IMAGE_NAME
	envVars := []string{
		"POSTGRES_USER=" + DevDbDefaultUser,
		"POSTGRES_PASSWORD=" + DevDbDefaultPass,
	}
	port := "5432"

	dh.StartContainerUI(imageName, envVars, port, "arkdb", nil)
}
