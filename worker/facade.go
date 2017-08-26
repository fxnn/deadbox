package worker

import (
	"fmt"
	"net/url"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
)

type facade struct {
	daemon.Daemon
	name                        string
	db                          *bolt.DB
	drop                        model.Drop
	dropUrl                     *url.URL
	updateRegistrationInterval  time.Duration
	registrationTimeoutDuration time.Duration
	registrator
}

func New(c config.Worker, db *bolt.DB) daemon.Daemon {
	drop := rest.NewClient(c.DropUrl)
	f := &facade{
		db:                         db,
		dropUrl:                    c.DropUrl,
		updateRegistrationInterval: time.Duration(c.UpdateRegistrationIntervalInSeconds) * time.Second,
		registrator: registrator{
			id:   generateWorkerId(),
			drop: drop,
			name: c.Name,
			registrationTimeoutDuration: time.Duration(c.RegistrationTimeoutInSeconds) * time.Second,
		},
	}
	f.Daemon = daemon.New(f.main)
	return f
}

func (f *facade) main(stop <-chan struct{}) error {
	if err := f.updateRegistration(); err != nil {
		err = fmt.Errorf("worker %s at drop %s could not be registered: %s", f.quotedNameAndId(), f.dropUrl, err)
		return err
	}

	updateRegistrationTicker := time.NewTicker(f.updateRegistrationInterval)
	defer updateRegistrationTicker.Stop()

	for {
		select {
		case <-updateRegistrationTicker.C:
			if err := f.updateRegistration(); err != nil {
				return fmt.Errorf("worker %s at drop %s could not update its registration: %s",
					f.quotedNameAndId(), f.dropUrl, err)
			}
		case <-stop:
			return nil
		}
	}
}
