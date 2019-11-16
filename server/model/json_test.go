package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWorker2Json(t *testing.T) {
	var err error
	var expectedId WorkerId = "id"
	var expectedTimeout time.Time = time.Now()

	bytes, err := json.Marshal(&Worker{
		Id:      expectedId,
		Timeout: expectedTimeout,
	})
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
		t.Fatalf("expected id %v, but got %v",
			expectedId, worker.Id)
	}
	if worker.Timeout.Sub(expectedTimeout) > 1*time.Millisecond || expectedTimeout.Sub(worker.Timeout) > 1*time.Millisecond {
		t.Fatalf("expected timeout %v, but got %v", expectedTimeout, worker.Timeout)
	}
}
