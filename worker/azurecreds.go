package worker

import (
	"log"
	"os"

	shell "github.com/katasec/utils/shell"
)

// getAzureCredsFromEnv reads Azure creds from env vars and sets them in pulumi config
func (w *Worker) getAzureCredsFromEnv() {

	// Define env vars to read
	log.Println("Reading Azure creds from env vars")
	envvars := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_TENANT_ID",
		"ARM_SUBSCRIPTION_ID",
	}

	// Check if all env vars are set
	ok := false
	for _, envvar := range envvars {
		if os.Getenv(envvar) != "" {
			ok = true
			shell.ExecShellCmd("pulumi config set azure-native:clientId " + os.Getenv(envvar))
		} else {
			ok = false
			log.Println("Env var " + envvar + " not set")
		}
	}

	// Exit if any env var is not set
	if !ok {
		log.Println("Some env vars not set, exitting...")
		os.Exit(1)
	}
}
