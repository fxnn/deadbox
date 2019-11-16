package it

import (
	"testing"
	"time"
)

func TestRegistration(t *testing.T) {

	daemon, drop := runDropDaemon(t)
	defer stopDaemon(daemon, t)

	worker, _ := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	// HINT: Give the worker some time to register
	time.Sleep(interactionSleepTime)

	actualWorkers, err := drop.Workers()
	if err != nil {
		t.Fatalf("could not read drop's worker list: %s", err)
	}

	assertNumberOfWorkers(actualWorkers, 1, t)
	actualWorker := actualWorkers[0]
	assertWorkerName(actualWorker, workerName, t)
	assertWorkerTimeoutInFuture(actualWorker, t)

}
