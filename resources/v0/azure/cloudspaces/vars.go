package cloudspaces

var (
	DefaultOctet1 = 10
	DefaultOctet2 = 16

	AllOctet2 = []int{}

	VnetPrefix = "vnet"
)

func init() {

	// Initialize Octet 2 range to  = 17..200
	for i := 17; i <= 200; i++ {
		AllOctet2 = append(AllOctet2, i)
	}
}
