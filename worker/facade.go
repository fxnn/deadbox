package worker

import (
	"crypto/rsa"
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
	privateKey                  *rsa.PrivateKey
}

func New(c *config.Worker, id string, db *bolt.DB, privateKey *rsa.PrivateKey, fingerprintLength uint, fingerprintChallengeLevel uint) Daemonized {
	drop := rest.NewClient(c.DropUrl, verifyByFingerprint(c.DropFingerprint, fingerprintLength, fingerprintChallengeLevel))
	f := &facade{
		db:                         db,
		dropUrl:                    c.DropUrl,
		privateKey:                 privateKey,
		updateRegistrationInterval: time.Duration(c.UpdateRegistrationIntervalInSeconds) * time.Second,
		registrations: registrations{
			id:   model.WorkerId(id),
			drop: drop,
			name: c.Name,
			registrationTimeoutDuration: time.Duration(c.RegistrationTimeoutInSeconds) * time.Second,
		},
		requests: requests{
			id:   model.WorkerId(id),
			drop: drop,
		},
		requestProcessors: &requestProcessors{
			processorsById: createRequestProcessorsByIdMap(c),
		},
	}
	f.Daemon = daemon.New(f.main)
	return f
}

func verifyByFingerprint(fingerprint string, fingerprintLength uint, fingerprintChallengeLevel uint) (verifyByFingerprint *crypto.VerifyByFingerprint) {
	if fingerprint != "" {
		verifyByFingerprint = &crypto.VerifyByFingerprint{
			Fingerprint:               fingerprint,
			FingerprintLength:         fingerprintLength,
			FingerprintChallengeLevel: fingerprintChallengeLevel,
		}
	}

	return
}

func (f *facade) main(stop <-chan struct{}) error {
	var err error

	publicKeyBytes, err := crypto.GeneratePublicKeyBytes(f.privateKey)
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
			if err := f.pollRequests(f.requestProcessors, f.privateKey); err != nil {
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
