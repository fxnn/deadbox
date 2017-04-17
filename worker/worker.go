package worker

import (
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/gone/log"
	"net/http"
	"net/url"
	"time"
)

type facade struct {
	db      *bolt.DB
	dropUrl *url.URL
}

func New(c config.Worker, db *bolt.DB) func() error {
	return facade{db, c.DropUrl}.Run
}

func (f *facade) Run() error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	select {
	case ticker.C:
		f.connectToDrop()
	}

	return nil
}
func (f *facade) connectToDrop() {
	resp, err := http.Get(f.dropUrl.ResolveReference("/worker").String())
	if err != nil {
		log.Warnln("couldn't connect to drop at ", f.dropUrl, ": ", err)
		return
	}
}
