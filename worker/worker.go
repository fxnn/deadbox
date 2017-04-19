package worker

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
	"net/url"
	"time"
)

const updateRegistrationInterval = 10 * time.Second
const registrationTimeoutDuration = 10 * updateRegistrationInterval

type facade struct {
	id      model.WorkerId
	db      *bolt.DB
	drop    model.Drop
	dropUrl *url.URL
}

func New(c config.Worker, db *bolt.DB) func() error {
	w := &facade{model.WorkerId(c.Name), db, rest.NewClient(c.DropUrl), c.DropUrl}
	return w.Run
}

func (f *facade) Run() error {
	if err := f.updateRegistration(); err != nil {
		err = fmt.Errorf("could not register worker %s at drop %s: %s", f.id, f.dropUrl, err)
		return err
	}

	updateRegistrationTicker := time.NewTicker(updateRegistrationInterval)
	defer updateRegistrationTicker.Stop()

	for {
		select {
		case <-updateRegistrationTicker.C:
			if err := f.updateRegistration(); err != nil {
				err = fmt.Errorf("could not update registration of worker %s at drop %s: %s", f.id, f.dropUrl, err)
				return err
			}
		}
	}
}
func (f *facade) updateRegistration() error {
	w := &model.Worker{Id: f.id, Timeout: time.Now().Add(registrationTimeoutDuration)}
	return f.drop.PutWorker(w)
}
