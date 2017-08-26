package worker

import (
	"fmt"
	"net/url"
	"time"

	"log"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/model"
	"github.com/fxnn/deadbox/rest"
)

const pollRequestInterval = 1 * time.Second

type Daemonized interface {
	daemon.Daemon
	Id() model.WorkerId
	Name() string
	QuotedNameAndId() string
}

type facade struct {
	daemon.Daemon
	registrator
	processor
	name                        string
	db                          *bolt.DB
	drop                        model.Drop
	dropUrl                     *url.URL
	updateRegistrationInterval  time.Duration
	registrationTimeoutDuration time.Duration
}

func New(c config.Worker, db *bolt.DB) Daemonized {
	drop := rest.NewClient(c.DropUrl)
	id := generateWorkerId()
	f := &facade{
		db:                         db,
		dropUrl:                    c.DropUrl,
		updateRegistrationInterval: time.Duration(c.UpdateRegistrationIntervalInSeconds) * time.Second,
		registrator: registrator{
			id:   id,
			drop: drop,
			name: c.Name,
			registrationTimeoutDuration: time.Duration(c.RegistrationTimeoutInSeconds) * time.Second,
		},
		processor: processor{
			id:   id,
			drop: drop,
		},
	}
	f.Daemon = daemon.New(f.main)
	return f
}

func (f *facade) main(stop <-chan struct{}) error {
	if err := f.updateRegistration(); err != nil {
		err = fmt.Errorf("worker %s at drop %s could not be registered: %s", f.QuotedNameAndId(), f.dropUrl, err)
		return err
	}

	updateRegistrationTicker := time.NewTicker(f.updateRegistrationInterval)
	defer updateRegistrationTicker.Stop()

	pollRequestTicker := time.NewTicker(pollRequestInterval)
	defer pollRequestTicker.Stop()

	for {
		select {
		case <-pollRequestTicker.C:
			// @todo #3 Replace pull with push mechanism (e.g. websocket)
			if err := f.pollRequests(); err != nil {
				log.Printf("worker %s at drop %s could not poll requests: %s", f.QuotedNameAndId(),
					f.dropUrl, err)
			}
		case <-updateRegistrationTicker.C:
			if err := f.updateRegistration(); err != nil {
				return fmt.Errorf("worker %s at drop %s could not update its registration: %s",
					f.QuotedNameAndId(), f.dropUrl, err)
			}
		case <-stop:
			return nil
		}
	}
}
