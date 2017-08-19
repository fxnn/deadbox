package it

import (
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/worker"
	"net/url"
	"os"
	"testing"
)

const workerDbFileName = "worker.boltdb"
const workerName = "itWorker"

func TestRegistration(t *testing.T) {

	var (
		workerCfg    config.Worker
		workerDb     *bolt.DB
		workerDaemon daemon.Daemon
		err          error
	)

	defer os.Remove(workerDbFileName)

	workerDb, err = bolt.Open(workerDbFileName, 0664, bolt.DefaultOptions)
	defer workerDb.Close()

	workerCfg = config.Worker{Name: workerName, DropUrl: parseUrlOrPanic("http://localhost:54123")}
	if err != nil {
		t.Fatalf("could not open Worker's BoltDB: %s", err)
	}

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
