package main

import (
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/drop"
)

func main() {
	cfg := config.Dummy()
	for _, dp := range cfg.Drops {
		drop.NewServer(dp.ListenAddress)
	}
}
