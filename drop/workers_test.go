package drop

import (
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/model"
	"os"
	"testing"
	"time"
)

const testDbFilename = "./test.boltdb"

func TestPutAndGet(t *testing.T) {

	var db *bolt.DB
	var sut *workers

	db = openTestDb()
	defer closeTestDb(db)

	sut = newWorkers(db)

	var given model.Worker = model.Worker{"id", time.Now()}
	sut.PutWorker(&given)

	var results []model.Worker = sut.Workers()
	if len(results) != 1 {
		t.Fatalf("expected response with 1 arg, but got %v",
			results)
	}
	var result model.Worker = results[0]
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
