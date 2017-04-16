package rest

import (
	"github.com/fxnn/deadbox/model"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestGetAllWorkers(t *testing.T) {
	rw := httptest.NewRecorder()
	rq := httptest.NewRequest(
		"GET",
		"http://localhost/worker",
		nil)

	r := newRouter(&mockDrop{})
	r.ServeHTTP(rw, rq)

	rp := rw.Result()
	body, _ := ioutil.ReadAll(rp.Body)

	if rp.StatusCode != 200 {
		t.Error("Status 200 expected, was", rp.StatusCode)
	}
	if string(body) != "[]\n" {
		t.Error("Empty JSON array expected, was ", string(body))
	}
	if rp.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type JSON expected, was",
			rp.Header.Get("Content-Type"))
	}
}

type mockDrop struct{}

func (*mockDrop) Workers() []model.Worker {
	return make([]model.Worker, 0)
}

func (*mockDrop) PutWorker(*model.Worker) {
	// no-op
}

func (*mockDrop) WorkerRequests(model.WorkerId) []model.WorkerRequest {
	return make([]model.WorkerRequest, 0)
}

func (*mockDrop) PutWorkerRequest(*model.WorkerRequest) {
	// no-op
}

func (*mockDrop) WorkerResponse(model.WorkerRequestId) []model.WorkerResponse {
	return make([]model.WorkerResponse, 0)
}

func (*mockDrop) PutWorkerResponse(*model.WorkerResponse) {
	// no-op
}
