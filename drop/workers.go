package drop

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/gone/log"
)

const workerBucketName = "worker"

type workers struct {
	db               *bolt.DB
	maxWorkerTimeout time.Duration
}

func (w *workers) Workers() ([]model.Worker, error) {
	var result []model.Worker
	var err error

	err = w.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(workerBucketName))
		if b == nil {
			log.Debugln("no worker bucket in db")
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var worker model.Worker
			if err := json.Unmarshal(v, &worker); err != nil {
				return err
			}
			result = append(result, worker)
		}

		return nil
	})
	return result, err
}

func (w *workers) PutWorker(worker *model.Worker) error {
	if worker.Id == "" {
		return fmt.Errorf("worker ID must not be empty")
	}
	if worker.Name == "" {
		return fmt.Errorf("worker Name must not be empty")
	}
	if worker.Timeout.Before(time.Now()) {
		return fmt.Errorf("worker timeout must be in the future")
	}
	if worker.Timeout.Sub(time.Now()) > w.maxWorkerTimeout {
		return fmt.Errorf("worker timeout cannot be more than %s in the future", w.maxWorkerTimeout.String())
	}

	return w.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(workerBucketName))
		if err != nil {
			return fmt.Errorf("couldn't create bucket %s: %v",
				workerBucketName, err)
		}

		v, err := json.Marshal(worker)
		if err != nil {
			return fmt.Errorf("couldn't marshal worker: %v", err)
		}
		if v == nil {
			return fmt.Errorf("unexpected nil for key %v", string(worker.Id))
		}

		err = b.Put([]byte(worker.Id), v)
		if err != nil {
			return fmt.Errorf("couldn't store worker %v: %v", v, err)
		}
		return nil
	})
}
