package dbtest

import (
	"fmt"
)

func init() {
	fmt.Println("This is init")
}

func Start() {
	fmt.Println("This is start!")

	acs := genCloudSpace()
	fmt.Println(acs)
}
