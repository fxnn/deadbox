package drop

import "github.com/fxnn/deadbox/model"
import (
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
)

// facade contains the implementation of model.Drop.
// As a facade, it redirects the method calls to the actual implementing
// structs.
type facade struct {
	name    string
	workers *workers
}

func New(c config.Drop, db *bolt.DB) model.Drop {
	return &facade{
		name:    c.Name,
		workers: newWorkers(db),
	}
}

func (f *facade) Workers() ([]model.Worker, error) {
	return f.workers.Workers()
}

func (f *facade) PutWorker(w *model.Worker) error {
	return f.workers.PutWorker(w)
}

func (*facade) WorkerRequests(model.WorkerId) []model.WorkerRequest {
	panic("implement me")
}

func (*facade) PutWorkerRequest(*model.WorkerRequest) {
	panic("implement me")
}

func (*facade) WorkerResponse(model.WorkerRequestId) []model.WorkerResponse {
	panic("implement me")
}

func (*facade) PutWorkerResponse(*model.WorkerResponse) {
	panic("implement me")
}
