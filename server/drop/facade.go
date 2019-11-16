package drop

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/server/config"
	"github.com/fxnn/deadbox/server/daemon"
	"github.com/fxnn/deadbox/server/model"
	"github.com/fxnn/deadbox/server/rest"
)

type Daemonized interface {
	model.Drop
	daemon.Daemon
}

// facade contains the implementation of model.Drop.
// As a facade, it redirects the method calls to the actual implementing
// structs.
type facade struct {
	daemon.Daemon
	name          string
	listenAddress string
	tls           rest.TLS
	*workers
	*requests
	*responses
}

func New(c *config.Drop, db *bolt.DB, tls rest.TLS) Daemonized {
	f := &facade{
		name:          c.Name,
		listenAddress: c.ListenAddress,
		tls:           tls,
		workers:       &workers{db, time.Duration(c.MaxWorkerTimeoutInSeconds) * time.Second},
		requests:      &requests{db, time.Duration(c.MaxRequestTimeoutInSeconds) * time.Second},
		responses:     &responses{db, time.Duration(c.MaxRequestTimeoutInSeconds) * time.Second},
	}
	f.Daemon = daemon.New(f.main)
	return f
}

func (f *facade) main(stop <-chan struct{}) error {
	server := rest.NewServer(f.listenAddress, f.tls, f)
	if err := server.StartServing(); err != nil {
		return fmt.Errorf("drop %s on %s could not be started: %s", f.quotedName(), f.listenAddress, err)
	}

	// @todo #10 secure drop against DoS and bruteforce attacks
	log.Println("drop", f.quotedName(), "on", f.listenAddress, "is now listening")
	for {
		select {
		case <-stop:
			log.Println("drop", f.quotedName(), "on", f.listenAddress, "shutting down")
			return server.Close()
		}
	}
}

func (f *facade) quotedName() string {
	return "'" + f.name + "'"
}
