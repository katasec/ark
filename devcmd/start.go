package devcmd

func Start() {
	d.RefreshConfig()

	setupDb()
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
