package it

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/drop"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
	"github.com/fxnn/deadbox/worker"
)

const workerDbFileName = "worker.boltdb"
const workerName = "itWorker"
const dropDbFileName = "drop.boltdb"
const dropName = "itDrop"
const port = "54123"

func assertWorkerTimeoutInFuture(actualWorker model.Worker, t *testing.T) {
	t.Helper()
	if actualWorker.Timeout.Before(time.Now()) {
		t.Fatalf("expected worker timeout to be in the future, but was %s", actualWorker.Timeout)
	}
}
func assertWorkerName(actualWorker model.Worker, workerName string, t *testing.T) {
	t.Helper()
	if string(actualWorker.Name) != workerName {
		t.Fatalf("expected worker to be %s, but was %v", workerName, actualWorker)
	}
}
func assertNumberOfWorkers(actualWorkers []model.Worker, expectedNumber int, t *testing.T) {
	t.Helper()
	if len(actualWorkers) != expectedNumber {
		t.Fatalf("expected %d workers, but got %v", expectedNumber, actualWorkers)
	}
}
func assertNumberOfRequests(actualRequests []model.WorkerRequest, expectedNumber int, t *testing.T) {
	t.Helper()
	if len(actualRequests) != expectedNumber {
		t.Fatalf("expected %d requests, but got %v", expectedNumber, actualRequests)
	}
}
func assertRequestId(actualRequest model.WorkerRequest, expectedId string, t *testing.T) {
	t.Helper()
	if string(actualRequest.Id) != expectedId {
		t.Fatalf("expected request to have id %s, but got %s", expectedId, actualRequest.Id)
	}
}

func assertResponseContentType(actualResponse model.WorkerResponse, expectedContentType string, t *testing.T) {
	t.Helper()
	if string(actualResponse.ContentType) != expectedContentType {
		t.Fatalf("expected response to have content type '%s', but got '%s' and content '%s'", expectedContentType, actualResponse.ContentType,
			actualResponse.Content)
	}
}
func assertResponseContent(actualResponse model.WorkerResponse, expectedContent string, t *testing.T) {
	t.Helper()
	if string(actualResponse.Content) != expectedContent {
		t.Fatalf("expected response to have content '%s', but got '%s'", expectedContent, actualResponse.Content)
	}
}

func runDropDaemon(t *testing.T) (daemon.Daemon, model.Drop) {
	t.Helper()

	cfg := config.Drop{
		Name:                       dropName,
		ListenAddress:              ":" + port,
		MaxRequestTimeoutInSeconds: config.DefaultMaxRequestTimeoutInSeconds,
		MaxWorkerTimeoutInSeconds:  config.DefaultMaxWorkerTimeoutInSeconds,
	}
	db, err := bolt.Open(dropDbFileName, 0664, bolt.DefaultOptions)
	if err != nil {
		t.Fatalf("could not open Drop's BoltDB: %s", err)
	}

	dropDaemon := drop.New(cfg, db)
	dropDaemon.OnStop(func() error {
		if err := db.Close(); err != nil {
			return err
		}
		if err := os.Remove(dropDbFileName); err != nil {
			return err
		}
		return nil
	})
	dropDaemon.Start()

	dropClient := rest.NewClient(parseUrlOrPanic("http://localhost:" + port))

	return dropDaemon, dropClient
}

func runWorkerDaemon(t *testing.T) worker.Daemonized {
	t.Helper()

	cfg := config.Worker{
		Name:    workerName,
		DropUrl: parseUrlOrPanic("http://localhost:" + port),
		RegistrationTimeoutInSeconds:        config.DefaultRegistrationTimeoutInSeconds,
		UpdateRegistrationIntervalInSeconds: config.DefaultUpdateRegistrationIntervalInSeconds,
	}
	db, err := bolt.Open(workerDbFileName, 0664, bolt.DefaultOptions)
	if err != nil {
		t.Fatalf("could not open Worker's BoltDB: %s", err)
	}

	workerDaemon := worker.New(cfg, db)
	workerDaemon.OnStop(func() error {
		if err := db.Close(); err != nil {
			return err
		}
		if err := os.Remove(workerDbFileName); err != nil {
			return err
		}
		return nil
	})
	workerDaemon.Start()

	return workerDaemon
}

func stopDaemon(d daemon.Daemon, t *testing.T) {
	t.Helper()
	err := d.Stop()
	if err != nil {
		t.Error(err)
	}
}

func parseUrlOrPanic(s string) *url.URL {
	result, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return result
}
