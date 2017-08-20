package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/drop"
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
	daemons := startDaemons(cfg)

	waitForShutdownRequest()

	log.Println("Shutting down gracefully")
	shutdownDaemons(daemons)
}

func shutdownDaemons(daemons []daemon.Daemon) {
	// TODO: Graceful HTTP shutdown with Go1.8
	for _, d := range daemons {
		if err := d.Stop(); err != nil {
			log.Println(err)
		}
	}
}

func startDaemons(cfg *config.Application) []daemon.Daemon {
	var daemons []daemon.Daemon = make([]daemon.Daemon, 0, len(cfg.Drops)+len(cfg.Workers))

	for _, dp := range cfg.Drops {
		daemons = append(daemons, serveDrop(dp, cfg))
	}

	for _, wk := range cfg.Workers {
		daemons = append(daemons, runWorker(wk, cfg))
	}

	return daemons
}
func waitForShutdownRequest() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

func runWorker(wcfg config.Worker, acfg *config.Application) daemon.Daemon {
	var b *bolt.DB = openDb(acfg, wcfg.Name)
	var d daemon.Daemon = worker.New(wcfg, b)
	d.OnStop(b.Close)
	d.Start()

	return d
}

func serveDrop(dcfg config.Drop, acfg *config.Application) daemon.Daemon {
	var b *bolt.DB = openDb(acfg, dcfg.Name)
	var d daemon.Daemon = drop.New(dcfg, b)
	d.OnStop(b.Close)
	d.Start()

	return d
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
