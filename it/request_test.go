package it

import (
	"github.com/fxnn/deadbox/model"
	"testing"
	"time"
)

const workerRequestId = "workerRequestId"

func TestRequest(t *testing.T) {
	t.Logf("TestRequest")

	drop := runDropDaemon(t)
	defer stopDaemon(drop, t)

	worker := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	var (
		err      error
		requests []model.WorkerRequest
		request  model.WorkerRequest = model.WorkerRequest{
			Id:      workerRequestId,
			Timeout: time.Now().Add(10 * time.Second),
		}
	)

	if err := drop.PutWorkerRequest(workerName, &request); err != nil {
		t.Fatalf("filing a new request failed: %s", err)
	}
	if requests, err = drop.WorkerRequests(workerName); err != nil {
		t.Fatalf("receiving a previously filed request failed: %s", err)
	}

	assertNumberOfRequests(requests, 1, t)
	actualRequest := requests[0]
	assertRequestId(actualRequest, workerRequestId, t)
	// TODO: Verify Id, Timeout and Content

}

func TestDuplicateRequestFails(t *testing.T) {
	t.Logf("TestRequest")

	drop := runDropDaemon(t)
	defer stopDaemon(drop, t)

	worker := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	var (
		request model.WorkerRequest = model.WorkerRequest{
			Id:      workerRequestId,
			Timeout: time.Now().Add(10 * time.Second),
		}
	)

	if err := drop.PutWorkerRequest(workerName, &request); err != nil {
		t.Fatalf("filing a new request failed: %s", err)
	}
	if err := drop.PutWorkerRequest(workerName, &request); err == nil {
		t.Fatalf("filing a request twice should've failed")
	}

}
