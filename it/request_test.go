package it

import (
	"testing"
	"time"

	"encoding/json"
	"fmt"

	"github.com/fxnn/deadbox/crypto"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/request/echo"
)

const (
	workerRequestId = "workerRequestId"
)

func TestRequest(t *testing.T) {

	daemon, drop := runDropDaemon(t)
	defer stopDaemon(daemon, t)

	worker, workerKey := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	// HINT: drop and worker some time to settle
	time.Sleep(interactionSleepTime)

	var (
		err      error
		requests []model.WorkerRequest
		request  = *encrypted(workerKey, t, &model.WorkerRequest{
			Id:          workerRequestId,
			Timeout:     time.Now().Add(10 * time.Second),
			ContentType: "application/json",
			Content:     echoRequest("test content"),
		})
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
	assertRequestEncryptionType(actualRequest, "encryptionType:github.com/fxnn/deadbox:AESPlusRSA:1.0", t)
	assertRequestContentContains(actualRequest, ":::", t)
	// @todo #2 Verify Id, Timeout and Content

	// HINT: Give worker time to send response
	time.Sleep(1000 * time.Millisecond)

	if response, err = drop.WorkerResponse(worker.Id(), request.Id); err != nil {
		t.Fatalf("receiving the response failed: %s", err)
	}
	assertResponseContentType(response, "application/json", t)
	assertResponseContent(response, "{\"echo\":\"test content\",\"requestProcessorId\":\"request-processor:github.com/fxnn/deadbox:echo:1.0\"}", t)

}

func encrypted(publicKeyBytes []byte, t *testing.T, request *model.WorkerRequest) *model.WorkerRequest {
	t.Helper()

	publicKey, err := crypto.UnmarshalPublicKey(publicKeyBytes)
	if err != nil {
		t.Fatalf("failed to unmarshal public key: %s", err)
	}

	contentEncrypted, encryptionType, err := crypto.EncryptRequest(request.Content, publicKey)
	if err != nil {
		t.Fatalf("failed to encrypt content: %s", err)
	}

	request.Content = contentEncrypted
	request.EncryptionType = encryptionType

	return request
}

func echoRequest(echoString string) []byte {
	var (
		m   = make(map[string]string)
		b   []byte
		err error
	)
	m["requestProcessorId"] = echo.RequestProcessorId
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

	worker, _ := runWorkerDaemon(t)
	defer stopDaemon(worker, t)

	// HINT: drop and worker some time to settle
	time.Sleep(interactionSleepTime)

	var request = model.WorkerRequest{
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
