package dbtest

import (
	"fmt"
)

func init() {
	// Register the driver.
	//	sql.Register("dbtest", &Driver{})
	fmt.Println("This is init")
}

func Start() {
	fmt.Println("This is start!")

	acs := genCloudSpace()
	fmt.Println(acs)
}
