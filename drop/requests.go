package drop

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/gone/log"
	"time"
)

const requestBucketName = "request"

type requests struct {
	db *bolt.DB
}

func newRequests(db *bolt.DB) *requests {
	return &requests{db}
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

func (w *requests) PutWorkerRequest(id model.WorkerId, request *model.WorkerRequest) error {
	// TODO: Validate name, and also the rest of the object
	if request.Id == "" {
		return fmt.Errorf("request ID must not be empty")
	}
	if request.Timeout.Before(time.Now()) {
		return fmt.Errorf("request timeout must be in the future")
	}

	return w.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(requestBucketName))
		if err != nil {
			return fmt.Errorf("bucket '%s' could not be created: %v", requestBucketName, err)
		}

		wb, err := b.CreateBucketIfNotExists([]byte(id))
		if err != nil {
			return fmt.Errorf("worker bucket '%s' could not be created: %v", id, err)
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
