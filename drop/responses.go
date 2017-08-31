package drop

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/gone/log"
)

type responses struct {
	db                 *bolt.DB
	maxResponseTimeout time.Duration
}

func (w *responses) WorkerResponse(id model.WorkerId, requestId model.WorkerRequestId) (result model.WorkerResponse, err error) {
	err = w.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(responseBucketName))
		if b == nil {
			log.Debugln("no response bucket in db")
			return nil
		}

		wb := b.Bucket([]byte(id))
		if wb == nil {
			log.Debugf("no response bucket for worker %s in db", id)
			return nil
		}

		v := wb.Get([]byte(requestId))
		if v == nil {
			return fmt.Errorf("no response in db for worker %s and requestId %s", id, requestId)
		}

		result = model.WorkerResponse{}
		if err := json.Unmarshal(v, &result); err != nil {
			return fmt.Errorf("response could not be unmarshalled from DB: %s", err)
		}

		return nil
	})

	return
}

func (w *responses) PutWorkerResponse(workerId model.WorkerId, requestId model.WorkerRequestId, response *model.WorkerResponse) error {
	if response.Timeout.Before(time.Now()) {
		return fmt.Errorf("response timeout must be in the future")
	}
	if response.Timeout.Sub(time.Now()) > w.maxResponseTimeout {
		return fmt.Errorf("response timeout cannot be more than %s in the future", w.maxResponseTimeout.String())
	}

	return w.db.Update(func(tx *bolt.Tx) error {
		if err := assertWorkerExists(tx, workerId); err != nil {
			return err
		}
		if err := deleteRequest(tx, workerId, requestId); err != nil {
			return err
		}
		return createResponse(tx, workerId, requestId, response)
	})
}
