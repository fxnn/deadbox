package drop

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/server/model"
)

type requests struct {
	db                *bolt.DB
	maxRequestTimeout time.Duration
}

func (w *requests) WorkerRequests(id model.WorkerId) ([]model.WorkerRequest, error) {
	var result []model.WorkerRequest
	var err error

	err = w.db.View(func(tx *bolt.Tx) error {
		// NOTE: this is a read-only TX, bucket expected to be created on worker registration
		if err := assertWorkerExists(tx, id); err != nil {
			return err
		}

		if wb, err := findOrCreateRequestBucket(tx, id); err != nil {
			return err
		} else {
			c := wb.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var request model.WorkerRequest
				if err = json.Unmarshal(v, &request); err != nil {
					return fmt.Errorf("request could not be unmarshalled from DB: %s", err)
				}
				result = append(result, request)
			}

			return nil
		}
	})
	return result, err
}

func (w *requests) PutWorkerRequest(workerId model.WorkerId, request *model.WorkerRequest) error {
	if request.Id == "" {
		return fmt.Errorf("request ID must not be empty")
	}
	if request.Timeout.Before(time.Now()) {
		return fmt.Errorf("request timeout must be in the future")
	}
	if request.Timeout.Sub(time.Now()) > w.maxRequestTimeout {
		return fmt.Errorf("request timeout cannot be more than %s in the future", w.maxRequestTimeout.String())
	}

	return w.db.Update(func(tx *bolt.Tx) error {
		if err := assertWorkerExists(tx, workerId); err != nil {
			return err
		}

		return createRequest(tx, workerId, request.Id, request)
	})
}
