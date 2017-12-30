package worker

import (
	"fmt"
	"net/url"
	"time"

	"log"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/crypto"
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
	registrations
	requests
	*requestProcessors
	name                        string
	db                          *bolt.DB
	drop                        model.Drop
	dropUrl                     *url.URL
	updateRegistrationInterval  time.Duration
	registrationTimeoutDuration time.Duration
	privateKeyBytes             []byte
}

func New(c config.Worker, db *bolt.DB, privateKeyBytes []byte) Daemonized {
	drop := rest.NewClient(c.DropUrl)
	id := generateWorkerId()
	f := &facade{
		db:                         db,
		dropUrl:                    c.DropUrl,
		privateKeyBytes:            privateKeyBytes,
		updateRegistrationInterval: time.Duration(c.UpdateRegistrationIntervalInSeconds) * time.Second,
		registrations: registrations{
			id:   id,
			drop: drop,
			name: c.Name,
			registrationTimeoutDuration: time.Duration(c.RegistrationTimeoutInSeconds) * time.Second,
		},
		requests: requests{
			id:   id,
			drop: drop,
		},
		requestProcessors: &requestProcessors{
			processorsById: createRequestProcessorsByIdMap(c),
		},
	}
	f.Daemon = daemon.New(f.main)
	return f
}

func (f *facade) main(stop <-chan struct{}) error {
	var err error

	privateKey, err := crypto.UnmarshalPrivateKeyFromPEMBytes(f.privateKeyBytes)
	if err != nil {
		return fmt.Errorf("worker %s could not read its private key from file %s: %s", f.QuotedNameAndId(), f.privateKeyBytes, err)
	}
	publicKeyBytes, err := crypto.GeneratePublicKeyBytes(privateKey)
	if err != nil {
		return fmt.Errorf("worker %s could not export its public key: %s", f.QuotedNameAndId(), err)
	}

	if err = f.updateRegistration(publicKeyBytes); err != nil {
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
			if err := f.pollRequests(f.requestProcessors, privateKey); err != nil {
				log.Printf("worker %s at drop %s could not poll requests: %s", f.QuotedNameAndId(),
					f.dropUrl, err)
			}
		case <-updateRegistrationTicker.C:
			if err := f.updateRegistration(publicKeyBytes); err != nil {
				return fmt.Errorf("worker %s at drop %s could not update its registration: %s",
					f.QuotedNameAndId(), f.dropUrl, err)
			}
		case <-stop:
			return nil
		}
	}
}

func GeneratePrivateKeyBytes() ([]byte, error) {
	if key, err := crypto.GeneratePrivateKey(); err != nil {
		return nil, fmt.Errorf("could not generate private key: %s", err)
	} else {
		return crypto.MarshalPrivateKeyToPEMBytes(key), nil
	}
}
