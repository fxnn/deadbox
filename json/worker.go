package json

import (
	"github.com/fxnn/deadbox/model"
	"time"
)

type Worker struct {
	IdVal      string
	TimeoutVal time.Time
}

func (w *Worker) Id() model.WorkerId {
	return model.WorkerId(w.IdVal)
}
func (w *Worker) Timeout() time.Time {
	return w.TimeoutVal
}

func AsWorkers(workers []model.Worker) []*Worker {
	result := make([]*Worker, len(workers))
	for i, worker := range workers {
		result[i] = AsWorker(worker)
	}
	return result
}
func AsWorker(worker model.Worker) *Worker {
	return &Worker{IdVal: string(worker.Id()), TimeoutVal: worker.Timeout()}
}
func NewWorker(id string, timeout time.Time) *Worker {
	return &Worker{IdVal: id, TimeoutVal: timeout}
}
