package drop

import (
	jsonenc "encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/json"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/gone/log"
)

const workerBucketName = "worker"

type workers struct {
	db *bolt.DB
}

func newWorkers(db *bolt.DB) *workers {
	return &workers{db}
}

func (w *workers) Workers() []model.Worker {
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
			var worker *json.Worker = new(json.Worker)
			if err := jsonenc.Unmarshal(v, worker); err != nil {
				return err
			}
			result = append(result, worker)
		}

		return nil
	})
	if err != nil {
		log.Debugln("error reading workers:", err)
	}
	return result
}

func (w *workers) PutWorker(worker model.Worker) {
	var err error
	// TODO: Validate name, and also the rest of the object

	err = w.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(workerBucketName))
		if err != nil {
			return fmt.Errorf("couldn't create bucket %s: %v",
				workerBucketName, err)
		}

		v, err := jsonenc.Marshal(json.AsWorker(worker))
		if err != nil {
			return fmt.Errorf("couldn't marshal worker: %v", err)
		}
		if v == nil {
			return fmt.Errorf("unexpected nil for key %v",
				worker.Id())
		}

		err = b.Put([]byte(worker.Id()), v)
		if err != nil {
			return fmt.Errorf("couldn't store worker %v: %v",
				v, err)
		}
		return nil
	})
	if err != nil {
		log.Debugln("error putting worker:", err)
		// TODO: Return error, so that we can return a 500
	}
}
