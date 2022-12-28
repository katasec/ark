package cli

// var (
// 	d = NewDevCmd()
// )

func Setup() {
	// Create Cloud resources with Pulumi
	d.Setup()

}

func SetupDelete() {
	// Delete Cloud resources with Pulumi
	d.Delete()
}
