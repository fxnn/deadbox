package drop

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/gone/log"
)

const requestBucketName = "request"

type requests struct {
	db                *bolt.DB
	maxRequestTimeout time.Duration
}

func (w *requests) WorkerRequests(id model.WorkerId) ([]model.WorkerRequest, error) {
	var result []model.WorkerRequest
	var err error

	err = w.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(requestBucketName))
		if b == nil {
			log.Debugln("no request bucket in db")
			return nil
		}

		wb := b.Bucket([]byte(id))
		if wb == nil {
			log.Debugf("no request bucket for worker %s in db", id)
			return nil
		}

		c := wb.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var request model.WorkerRequest
			if err := json.Unmarshal(v, &request); err != nil {
				return fmt.Errorf("request could not be unmarshalled from DB: %s", err)
			}
			result = append(result, request)
		}

		return nil
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
		b := tx.Bucket([]byte(workerBucketName))
		if b == nil {
			return fmt.Errorf("worker '%s' does not exist: no worker bucket in db", workerId)
		}
		if b.Get([]byte(workerId)) == nil {
			return fmt.Errorf("worker '%s' does not exist", workerId)
		}

		b, err := tx.CreateBucketIfNotExists([]byte(requestBucketName))
		if err != nil {
			return fmt.Errorf("bucket '%s' could not be created: %v", requestBucketName, err)
		}

		wb, err := b.CreateBucketIfNotExists([]byte(workerId))
		if err != nil {
			return fmt.Errorf("worker bucket '%s' could not be created: %v", workerId, err)
		}

		if wb.Get([]byte(request.Id)) != nil {
			return fmt.Errorf("request '%s' already filed", string(request.Id))
		}

		v, err := json.Marshal(request)
		if err != nil {
			return fmt.Errorf("request could not be marshalled: %v", err)
		}
		if v == nil {
			return fmt.Errorf("unexpected nil for key %v", string(request.Id))
		}

		err = wb.Put([]byte(request.Id), v)
		if err != nil {
			return fmt.Errorf("request could not be stored: %v", err)
		}
		return nil
	})
}
