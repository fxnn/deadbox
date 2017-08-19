package it

import (
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/drop"
	"github.com/fxnn/deadbox/worker"
	"net/url"
	"os"
	"testing"
)

const workerDbFileName = "worker.boltdb"
const workerName = "itWorker"
const dropDbFileName = "drop.boltdb"
const dropName = "itDrop"
const itPort = "54123"

func TestRegistration(t *testing.T) {

	var (
		dropCfg    config.Drop
		dropDb     *bolt.DB
		dropDaemon daemon.Daemon
		err        error
	)

	defer os.Remove(dropDbFileName)

	dropDb, err = bolt.Open(dropDbFileName, 0664, bolt.DefaultOptions)
	if err != nil {
		t.Fatalf("could not open Drop's BoltDB: %s", err)
	}
	defer dropDb.Close()

	dropCfg = config.Drop{Name: dropName, ListenAddress: ":" + itPort}
	dropDaemon = drop.New(dropCfg, dropDb)
	defer dropDaemon.Stop()

	dropDaemon.Start()

	var (
		workerCfg    config.Worker
		workerDb     *bolt.DB
		workerDaemon daemon.Daemon
	)

	defer os.Remove(workerDbFileName)

	workerDb, err = bolt.Open(workerDbFileName, 0664, bolt.DefaultOptions)
	if err != nil {
		t.Fatalf("could not open Worker's BoltDB: %s", err)
	}
	defer workerDb.Close()

	workerCfg = config.Worker{Name: workerName, DropUrl: parseUrlOrPanic("http://localhost:" + itPort)}
	workerDaemon = worker.New(workerCfg, workerDb)
	defer workerDaemon.Stop()

	workerDaemon.Start()

}

func parseUrlOrPanic(s string) *url.URL {
	result, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return result
}
