package rest

import (
	"github.com/fxnn/deadbox/model"
	"time"
)

type jsonWorker struct {
	id      string
	timeout time.Time
}

func (w *jsonWorker) Id() model.WorkerId {
	return model.WorkerId(w.id)
}
func (w *jsonWorker) Timeout() time.Time {
	return w.timeout
}

func asJsonWorker(worker model.Worker) jsonWorker {
	return jsonWorker{id: string(worker.Id()), timeout: worker.Timeout()}
}
func asJsonWorkers(workers []model.Worker) []jsonWorker {
	result := make([]jsonWorker, len(workers))
	for i, worker := range workers {
		result[i] = asJsonWorker(worker)
	}
	return result
}

type jsonWorkerRequest struct {
	id      string
	timeout time.Time
}

func (r *jsonWorkerRequest) Id() model.WorkerRequestId {
	return model.WorkerRequestId(r.id)
}
func (r *jsonWorkerRequest) Timeout() time.Time {
	return r.timeout
}

func asJsonWorkerRequest(request model.WorkerRequest) jsonWorkerRequest {
	return jsonWorkerRequest{id: string(request.Id()), timeout: request.Timeout()}
}
func asJsonWorkerRequests(requests []model.WorkerRequest) []jsonWorkerRequest {
	result := make([]jsonWorkerRequest, len(requests))
	for i, request := range requests {
		result[i] = asJsonWorkerRequest(request)
	}
	return result
}
