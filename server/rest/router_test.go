package rest

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/fxnn/deadbox/server/model"
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

func (*mockDrop) Workers() ([]model.Worker, error) {
	return make([]model.Worker, 0), nil
}

func (*mockDrop) PutWorker(*model.Worker) error {
	return nil
}

func (*mockDrop) WorkerRequests(model.WorkerId) ([]model.WorkerRequest, error) {
	return make([]model.WorkerRequest, 0), nil
}

func (*mockDrop) PutWorkerRequest(model.WorkerId, *model.WorkerRequest) error {
	return nil
}

func (*mockDrop) WorkerResponse(model.WorkerId, model.WorkerRequestId) (model.WorkerResponse, error) {
	return model.WorkerResponse{}, nil
}

func (*mockDrop) PutWorkerResponse(model.WorkerId, model.WorkerRequestId, *model.WorkerResponse) error {
	return nil
}
