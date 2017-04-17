package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/drop"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
	"github.com/fxnn/deadbox/worker"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const dbFileExtension = "boltdb"

func main() {
	var cfg *config.Application = config.Dummy()

	for _, dp := range cfg.Drops {
		go serveDrop(dp, cfg)
	}

	for _, wk := range cfg.Workers {
		go runWorker(wk, cfg)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Shutting down gracefully")
	// TODO: Graceful HTTP shutdown with Go1.8
}

func runWorker(wcfg config.Worker, acfg *config.Application) {
	var b *bolt.DB = openDb(acfg, wcfg.Name)
	defer closeDb(b)

	var runWorker func() error = worker.New(wcfg, b)
	log.Fatalln(runWorker())
}

func serveDrop(dcfg config.Drop, acfg *config.Application) {
	var b *bolt.DB = openDb(acfg, dcfg.Name)
	defer closeDb(b)

	var dp model.Drop = drop.New(dcfg, b)
	log.Println("Drop", dcfg.Name, "listening on", dcfg.ListenAddress)
	log.Fatalln(rest.NewServer(dcfg.ListenAddress, dp).Serve())
}

func closeDb(b *bolt.DB) {
	if err := b.Close(); err != nil {
		log.Fatal(err)
	}
}

func openDb(cfg *config.Application, name string) *bolt.DB {
	boltOptions := &bolt.Options{Timeout: 10 * time.Second}

	fileName := dbFileName(cfg, name)
	db, err := bolt.Open(fileName, 0660, boltOptions)
	if err != nil {
		panic(fmt.Errorf(
			"couldn't open bolt DB %s: %s",
			fileName, err,
		))
	}
	log.Println("Database opened:", fileName)

	return db
}

func dbFileName(cfg *config.Application, name string) string {
	return filepath.Join(cfg.DbPath, name+"."+dbFileExtension)
}
