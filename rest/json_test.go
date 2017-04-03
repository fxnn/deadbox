package rest

import (
	"github.com/fxnn/deadbox/model"
	"testing"
)

func TestJsonWorker_Interface(*testing.T) {
	var _ model.Worker = &jsonWorker{}
}

func TestJsonWorkerRequest_Interface(*testing.T) {
	var _ model.WorkerRequest = &jsonWorkerRequest{}
}
