package json

import (
	"github.com/fxnn/deadbox/model"
	"time"
)

type Worker struct {
	id      string
	timeout time.Time
}

func (w *Worker) Id() model.WorkerId {
	return model.WorkerId(w.id)
}
func (w *Worker) Timeout() time.Time {
	return w.timeout
}

func NewWorker(worker model.Worker) Worker {
	return Worker{id: string(worker.Id()), timeout: worker.Timeout()}
}
func NewWorkers(workers []model.Worker) []Worker {
	result := make([]Worker, len(workers))
	for i, worker := range workers {
		result[i] = NewWorker(worker)
	}
	return result
}
