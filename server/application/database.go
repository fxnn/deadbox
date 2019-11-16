package application

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/server/config"
)

func OpenDb(cfg *config.Application, name string) *bolt.DB {
	boltOptions := &bolt.Options{Timeout: 10 * time.Second}

	fileName := dbFileName(cfg, name)
	db, err := bolt.Open(fileName, filePermUserAndGroupCanReadOrWrite, boltOptions)
	if err != nil {
		panic(fmt.Errorf(
			"couldn't open bolt DB %s: %s",
			fileName, err,
		))
	}
	log.Println("Database opened:", fileName)

	return db
}
