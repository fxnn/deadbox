package main

import (
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/rest"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var cfg *config.Application = config.Dummy()

	for _, dp := range cfg.Drops {
		go func(drop config.Drop) {
			wg.Add(1)
			// FIXME: Create drop instance
			rest.NewServer(dp.ListenAddress, nil)
		}(dp)
	}

	wg.Wait()
}
