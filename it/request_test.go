package it

import (
	"github.com/fxnn/deadbox/model"
	"testing"
)

func TestRequest(t *testing.T) {
	t.Logf("TestRequest")

	drop := runDropDaemon(t)
	defer stopDaemon(drop, t)

	worker := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	var (
		err      error
		requests []model.WorkerRequest
	)

	if err := drop.PutWorkerRequest(workerName, &model.WorkerRequest{}); err != nil {
		t.Fatalf("filing a new request failed: %s", err)
	}
	if requests, err = drop.WorkerRequests(workerName); err != nil {
		t.Fatalf("receiving a previously filed request failed: %s", err)
	}

	assertNumberOfRequests(requests, 1, t)
	// TODO: Verify Id, Timeout and Content

}
