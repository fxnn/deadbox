package it

import (
	"testing"
	"time"
)

func TestRegistration(t *testing.T) {

	drop := runDropDaemon(t)
	defer drop.Stop()

	worker := runWorkerDaemon(t)
	defer worker.Stop()

	// HINT: Give the worker some time to register
	time.Sleep(200 * time.Millisecond)

	actualWorkers, err := drop.Workers()
	if err != nil {
		t.Fatalf("could not read drop's worker list: %s", err)
	}

	assertNumberOfWorkers(actualWorkers, 1, t)
	actualWorker := actualWorkers[0]
	assertWorkerName(actualWorker, workerName, t)
	assertWorkerTimeoutInFuture(actualWorker, t)

}
