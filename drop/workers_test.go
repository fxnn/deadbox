package drop

import (
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/model"
)

const testDbFilename = "./test.boltdb"

func TestPutAndGet(t *testing.T) {

	db := openTestDb()
	defer closeTestDb(db)

	sut := &workers{db, 10 * time.Minute}
	err := sut.PutWorker(&model.Worker{Id: "id", Timeout: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatalf("sut.PutWorker() returned error: %s", err)
	}

	results, err := sut.Workers()
	if err != nil {
		t.Fatalf("sut.Workers() returned error: %s", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected response with 1 arg, but got %v", results)
	}
	result := results[0]
	if result.Id != "id" {
		t.Fatalf("got id %v", string(result.Id))
	}
}

func closeTestDb(db *bolt.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
	os.Remove(testDbFilename)
}

func openTestDb() *bolt.DB {
	db, err := bolt.Open(testDbFilename, 0660, nil)
	if err != nil {
		panic(err)
	}

	return db
}
