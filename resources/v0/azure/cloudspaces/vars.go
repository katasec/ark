package cloudspaces

var (
	DefaultOctet1 = 10
	DefaultOctet2 = 16

	Octet2Range = []int{}
)

func init() {

	// Initialize Octet2Range
	for i := 17; i <= 200; i++ {
		Octet2Range = append(Octet2Range, i)
	}
}
