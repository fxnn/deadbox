package it

import (
	"testing"
	"time"

	"encoding/json"
	"fmt"

	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/processor/echo"
)

const workerRequestId = "workerRequestId"

func TestRequest(t *testing.T) {

	daemon, drop := runDropDaemon(t)
	defer stopDaemon(daemon, t)

	worker := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	// HINT: drop and worker some time to settle
	time.Sleep(100 * time.Millisecond)

	var (
		err      error
		requests []model.WorkerRequest
		request  model.WorkerRequest = model.WorkerRequest{
			Id:          workerRequestId,
			Timeout:     time.Now().Add(10 * time.Second),
			ContentType: "application/json",
			Content:     echoRequest("test content"),
		}
		response model.WorkerResponse
	)

	if err = drop.PutWorkerRequest(worker.Id(), &request); err != nil {
		t.Fatalf("filing a new request failed: %s", err)
	}
	if requests, err = drop.WorkerRequests(worker.Id()); err != nil {
		t.Fatalf("receiving a previously filed request failed: %s", err)
	}

	assertNumberOfRequests(requests, 1, t)
	actualRequest := requests[0]
	assertRequestId(actualRequest, workerRequestId, t)
	// @todo #2 Verify Id, Timeout and Content

	// HINT: Give worker time to send response
	time.Sleep(1000 * time.Millisecond)

	if response, err = drop.WorkerResponse(worker.Id(), request.Id); err != nil {
		t.Fatalf("receiving the response failed: %s", err)
	}
	assertResponseContentType(response, "application/json", t)
	assertResponseContent(response, "{\"echo\":\"test content\",\"requestType\":\"github.com/fxnn/deadbox/processor/echo 1.0\"}", t)

}

func echoRequest(echoString string) []byte {
	var (
		m   map[string]string = make(map[string]string)
		b   []byte
		err error
	)
	m["requestType"] = echo.RequestType
	m["echo"] = echoString
	b, err = json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("marshalling echo request content failed: %s", err))
	}
	return b
}

func TestDuplicateRequestFails(t *testing.T) {

	daemon, drop := runDropDaemon(t)
	defer stopDaemon(daemon, t)

	worker := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	// HINT: drop and worker some time to settle
	time.Sleep(100 * time.Millisecond)

	var request model.WorkerRequest = model.WorkerRequest{
		Id:      workerRequestId,
		Timeout: time.Now().Add(10 * time.Second),
	}

	if err := drop.PutWorkerRequest(worker.Id(), &request); err != nil {
		t.Fatalf("filing a new request failed: %s", err)
	}
	if err := drop.PutWorkerRequest(worker.Id(), &request); err == nil {
		t.Fatalf("filing a request twice should've failed")
	}

}
