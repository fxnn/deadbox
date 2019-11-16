package worker

import (
	"time"

	"fmt"

	"github.com/fxnn/deadbox/server/model"
)

type registrations struct {
	id                          model.WorkerId
	name                        string
	drop                        model.Drop
	registrationTimeoutDuration time.Duration
}

func (r *registrations) updateRegistration(publicKey []byte) error {
	w := &model.Worker{
		Id:        r.id,
		Name:      r.name,
		Timeout:   time.Now().Add(r.registrationTimeoutDuration),
		PublicKey: publicKey,
	}
	return r.drop.PutWorker(w)
}

func (r *registrations) Id() model.WorkerId {
	return r.id
}

func (r *registrations) Name() string {
	return r.name
}

func (r *registrations) QuotedNameAndId() string {
	return fmt.Sprintf("'%s' (%s)", r.name, r.id)
}
