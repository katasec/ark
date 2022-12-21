package dev

func Setup() {

	// Create Cloud resources with PuUlumi
	createLocal()

	// Update config file with links to new pulumi cloud resources.
	refreshConfig()
}

func SetupDelete() {
	deleteLocal()

}
