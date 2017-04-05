package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/drop"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cfg *config.Application = config.Dummy()

	var db *bolt.DB = openDb(cfg)
	defer db.Close()

	for _, dp := range cfg.Drops {
		go serveDrop(dp, db)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Shutting down gracefully")
	db.Close()
	// TODO: Graceful HTTP shutdown with Go1.8
}

func openDb(cfg *config.Application) *bolt.DB {
	boltOptions := &bolt.Options{Timeout: 10 * time.Second}

	db, err := bolt.Open(cfg.DbFile, 0660, boltOptions)
	if err != nil {
		panic(fmt.Errorf("couldn't open bolt DB: %s", err))
	}
	log.Println("Database opened")

	// TODO: Create all needed buckets, if not yet existing

	return db
}

func serveDrop(cfg config.Drop, db *bolt.DB) {
	var dp model.Drop = drop.New(cfg.Name, db)
	log.Println("Drop", cfg.Name, "listening on", cfg.ListenAddress)
	log.Fatalln(rest.NewServer(cfg.ListenAddress, dp).Serve())
}
