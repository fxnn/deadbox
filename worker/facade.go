package worker

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
	"net/url"
	"time"
)

const updateRegistrationInterval = 10 * time.Second
const registrationTimeoutDuration = 10 * updateRegistrationInterval

type facade struct {
	daemon.Daemon
	id      model.WorkerId
	db      *bolt.DB
	drop    model.Drop
	dropUrl *url.URL
}

func New(c config.Worker, db *bolt.DB) daemon.Daemon {
	f := &facade{nil, model.WorkerId(c.Name), db, rest.NewClient(c.DropUrl), c.DropUrl}
	f.Daemon = daemon.New(f.main)
	return f
}

func (f *facade) main(stop <-chan struct{}) error {
	if err := f.updateRegistration(); err != nil {
		err = fmt.Errorf("worker %s at drop %s could not be registered: %s", f.quotedId(), f.dropUrl, err)
		return err
	}

	updateRegistrationTicker := time.NewTicker(updateRegistrationInterval)
	defer updateRegistrationTicker.Stop()

	for {
		select {
		case <-updateRegistrationTicker.C:
			if err := f.updateRegistration(); err != nil {
				return fmt.Errorf("worker %s at drop %s could not update its registration: %s",
					f.quotedId(), f.dropUrl, err)
			}
		case <-stop:
			return nil
		}
	}
}

func (f *facade) quotedId() string {
	return "'" + string(f.id) + "'"
}

func (f *facade) updateRegistration() error {
	w := &model.Worker{Id: f.id, Timeout: time.Now().Add(registrationTimeoutDuration)}
	return f.drop.PutWorker(w)
}
