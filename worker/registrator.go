package worker

import (
	"time"

	"crypto/rand"
	"encoding/base64"

	"fmt"

	"github.com/fxnn/deadbox/model"
)

const idBytesEntropy = 64

type registrator struct {
	id                          model.WorkerId
	name                        string
	drop                        model.Drop
	registrationTimeoutDuration time.Duration
}

func (r *registrator) updateRegistration() error {
	w := &model.Worker{
		Id:      r.id,
		Name:    r.name,
		Timeout: time.Now().Add(r.registrationTimeoutDuration),
	}
	return r.drop.PutWorker(w)
}

func (r *registrator) Id() model.WorkerId {
	return r.id
}

func (r *registrator) Name() string {
	return r.name
}

func (r *registrator) QuotedNameAndId() string {
	return fmt.Sprintf("'%s' (%s)", r.name, r.id)
}

func generateWorkerId() model.WorkerId {

	rawBytes := make([]byte, idBytesEntropy)
	if _, err := rand.Read(rawBytes); err != nil {
		panic(fmt.Sprint("couldn't generate random bytes for worker id:", err))
	}

	encoded := base64.RawURLEncoding.EncodeToString(rawBytes)
	return model.WorkerId(encoded)

}
