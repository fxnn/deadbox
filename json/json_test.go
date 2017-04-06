package json

import (
	"github.com/fxnn/deadbox/model"
	"testing"
)

func TestJsonWorker_Interface(*testing.T) {
	var _ model.Worker = &Worker{}
}

func TestJsonWorkerRequest_Interface(*testing.T) {
	var _ model.WorkerRequest = &WorkerRequest{}
}
