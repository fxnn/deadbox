package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/drop"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var cfg *config.Application = config.Dummy()

	var db *bolt.DB = openDb(cfg)

	for _, dp := range cfg.Drops {
		wg.Add(1)
		go serveDrop(wg, dp, db)
	}

	wg.Wait()
}
func openDb(cfg *config.Application) *bolt.DB {
	boltOptions := &bolt.Options{Timeout: 10 * time.Second}

	db, err := bolt.Open(cfg.DbFile, 0660, boltOptions)
	if err != nil {
		panic(fmt.Errorf("couldn't open bolt DB: %s", err))
	}

	// TODO: Create all needed buckets, if not yet existing

	return db
}

func serveDrop(wg sync.WaitGroup, cfg config.Drop, db *bolt.DB) {
	var dp model.Drop = drop.New(db)
	rest.NewServer(cfg.ListenAddress, dp).Serve()
	wg.Done()
}
