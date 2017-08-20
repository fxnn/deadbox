package drop

import "github.com/fxnn/deadbox/model"
import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/rest"
	"log"
)

type DaemonizedDrop interface {
	model.Drop
	daemon.Daemon
}

// facade contains the implementation of model.Drop.
// As a facade, it redirects the method calls to the actual implementing
// structs.
type facade struct {
	daemon.Daemon
	name          string
	listenAddress string
	workers       *workers
}

func New(c config.Drop, db *bolt.DB) DaemonizedDrop {
	f := &facade{
		name:          c.Name,
		listenAddress: c.ListenAddress,
		workers:       newWorkers(db),
	}
	f.Daemon = daemon.New(f.main)
	return f
}

func (f *facade) main(stop <-chan struct{}) error {
	server := rest.NewServer(f.listenAddress, f)
	if err := server.StartServing(); err != nil {
		return fmt.Errorf("drop %s on %s could not be started: %s", f.quotedName(), f.listenAddress, err)
	}

	log.Println("drop", f.quotedName(), "on", f.listenAddress, "is now listening")
	for {
		select {
		case <-stop:
			log.Println("drop", f.quotedName(), "on", f.listenAddress, "shutting down")
			return server.Close()
		}
	}
}

func (f *facade) quotedName() string {
	return "'" + f.name + "'"
}

func (f *facade) Workers() ([]model.Worker, error) {
	return f.workers.Workers()
}

func (f *facade) PutWorker(w *model.Worker) error {
	return f.workers.PutWorker(w)
}

func (*facade) WorkerRequests(model.WorkerId) ([]model.WorkerRequest, error) {
	return nil, fmt.Errorf("implement me")
}

func (*facade) PutWorkerRequest(model.WorkerId, *model.WorkerRequest) error {
	return fmt.Errorf("implement me")
}

func (*facade) WorkerResponse(model.WorkerId, model.WorkerRequestId) (model.WorkerResponse, error) {
	return model.WorkerResponse{}, fmt.Errorf("implement me")
}

func (*facade) PutWorkerResponse(model.WorkerId, model.WorkerRequestId, *model.WorkerResponse) error {
	return fmt.Errorf("implement me")
}
