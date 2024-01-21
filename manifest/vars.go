package manifest

import "github.com/katasec/ark/config"

var (
	arkConfig *config.Config
)

func init() {
	arkConfig = config.ReadConfig()
}
