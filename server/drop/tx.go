package drop

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/server/model"
)

const responseBucketName = "request"
const requestBucketName = "request"

func assertWorkerExists(tx *bolt.Tx, workerId model.WorkerId) error {
	b := tx.Bucket([]byte(workerBucketName))
	if b == nil {
		return fmt.Errorf("worker '%s' does not exist: no worker bucket in db", workerId)
	}
	if b.Get([]byte(workerId)) == nil {
		return fmt.Errorf("worker '%s' does not exist", workerId)
	}
	return nil
}

func createRequest(tx *bolt.Tx, workerId model.WorkerId, requestId model.WorkerRequestId, request *model.WorkerRequest) error {
	key, value, err := marshalRequestKeyValue(requestId, request)
	if err != nil {
		return err
	}

	if wb, err := findOrCreateRequestBucket(tx, workerId); err != nil {
		return err
	} else if wb.Get(key) != nil {
		return fmt.Errorf("request for '%s' already filed", string(requestId))
	} else if err := wb.Put(key, value); err != nil {
		return fmt.Errorf("request could not be stored: %v", err)
	}

	return nil
}

func marshalRequestKeyValue(requestId model.WorkerRequestId, request *model.WorkerRequest) ([]byte, []byte, error) {
	key := []byte(requestId)
	if value, err := json.Marshal(request); err != nil {
		return nil, nil, fmt.Errorf("request could not be marshalled: %v", err)
	} else if value == nil {
		return nil, nil, fmt.Errorf("unexpected nil for key %v", string(requestId))
	} else {
		return key, value, nil
	}
}

func deleteRequest(tx *bolt.Tx, workerId model.WorkerId, requestId model.WorkerRequestId) error {
	key := []byte(requestId)

	if wb, err := findOrCreateRequestBucket(tx, workerId); err != nil {
		return err
	} else if err := wb.Delete(key); err != nil {
		return fmt.Errorf("request bucket '%s' could not be deleted: %s", workerId, err)
	}

	return nil
}

func findOrCreateRequestBucket(tx *bolt.Tx, workerId model.WorkerId) (*bolt.Bucket, error) {
	if b, err := findOrCreateBucket(tx, []byte(requestBucketName)); err != nil {
		return nil, err
	} else if wb, err := findOrCreateBucket(b, []byte(workerId)); err != nil {
		return nil, err
	} else {
		return wb, nil
	}
}

func createResponse(tx *bolt.Tx, workerId model.WorkerId, requestId model.WorkerRequestId, response *model.WorkerResponse) error {
	key, value, err := marshalResponseKeyValue(requestId, response)
	if err != nil {
		return err
	}

	if wb, err := findOrCreateResponseBucket(tx, workerId); err != nil {
		return err
	} else if wb.Get(key) != nil {
		return fmt.Errorf("response for '%s' already filed", string(requestId))
	} else if err := wb.Put(key, value); err != nil {
		return fmt.Errorf("response could not be stored: %v", err)
	}

	return nil
}

func marshalResponseKeyValue(requestId model.WorkerRequestId, response *model.WorkerResponse) ([]byte, []byte, error) {
	key := []byte(requestId)
	if value, err := json.Marshal(response); err != nil {
		return nil, nil, fmt.Errorf("response could not be marshalled: %v", err)
	} else if value == nil {
		return nil, nil, fmt.Errorf("unexpected nil for key %v", string(requestId))
	} else {
		return key, value, nil
	}
}

func findOrCreateResponseBucket(tx *bolt.Tx, workerId model.WorkerId) (*bolt.Bucket, error) {
	b, err := findOrCreateBucket(tx, []byte(responseBucketName))
	if err != nil {
		return nil, err
	}

	wb, err := findOrCreateBucket(b, []byte(workerId))
	if err != nil {
		return nil, err
	}

	return wb, nil
}

type txOrBucket interface {
	Bucket([]byte) *bolt.Bucket
	Writable() bool
	CreateBucket([]byte) (*bolt.Bucket, error)
}

func findOrCreateBucket(tx txOrBucket, name []byte) (*bolt.Bucket, error) {
	b := tx.Bucket(name)
	if b != nil {
		return b, nil
	}
	if !tx.Writable() {
		return nil, fmt.Errorf("bucket '%s' does not exist", requestBucketName)
	}
	if b, err := tx.CreateBucket(name); err != nil {
		return nil, fmt.Errorf("bucket '%s' could not be created: %v", requestBucketName, err)
	} else {
		return b, nil
	}
}
