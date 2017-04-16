package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWorkerMarshalAndUnmarshal(t *testing.T) {
	var err error
	var expectedId WorkerId = "id"

	bytes, err := json.Marshal(&Worker{Id: expectedId, Timeout: time.Now()})
	if err != nil {
		t.Fatal("couldn't marshal worker: ", err)
	}

	worker := new(Worker)
	err = json.Unmarshal(bytes, worker)
	if err != nil {
		t.Fatal("couldn't unmarshal worker \"", string(bytes), "\": ",
			err)
	}

	if worker.Id != expectedId {
		t.Fatalf("expected id %v, but got %v", expectedId, worker.Id)
	}
}
