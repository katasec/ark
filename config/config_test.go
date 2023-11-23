package config

import (
	"fmt"
	"testing"
)

func TestShowConfig(t *testing.T) {
	cfg := ReadConfig()
	fmt.Println("Command Queue:" + cfg.CmdQ)
	fmt.Println("Response Queue:" + cfg.RespQ)
}
