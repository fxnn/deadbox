package drop

import "github.com/fxnn/deadbox/model"
import "github.com/boltdb/bolt"

// facade contains the implementation of model.Drop.
// As a facade, it redirects the method calls to the actual implementing
// structs.
type facade struct {
	name string
	db   *bolt.DB
}

func New(name string, db *bolt.DB) model.Drop {
	return &facade{name, db}
}

func (*facade) Workers() []model.Worker {
	panic("implement me")
}

func (*facade) PutWorker(model.Worker) {
	panic("implement me")
}

func (*facade) WorkerRequests(model.WorkerId) []model.WorkerRequest {
	panic("implement me")
}

func (*facade) PutWorkerRequest(model.WorkerRequest) {
	panic("implement me")
}

func (*facade) WorkerResponse(model.WorkerRequestId) []model.WorkerResponse {
	panic("implement me")
}

func (*facade) PutWorkerResponse(model.WorkerResponse) {
	panic("implement me")
}
