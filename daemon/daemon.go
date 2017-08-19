package daemon

import (
	"fmt"
	"log"
)

type Daemon interface {
	Running() bool
	Start() error
	Stop() error
}

type daemon struct {
	main    func(stop <-chan struct{}) error
	stop    chan struct{}
	stopped chan struct{}
}

func New(main func(stop <-chan struct{}) error) Daemon {
	return &daemon{main, nil, nil}
}

func (d *daemon) Running() bool {
	return d.stop != nil
}

func (d *daemon) Start() error {
	// TODO: Make thread safe
	if !d.Running() {
		if d.main == nil {
			return fmt.Errorf("daemon could not be started: %s", "no main set")
		}

		d.stop = make(chan struct{})
		d.stopped = make(chan struct{})
		go func() {
			defer d.closeStopChannel()
			defer d.closeStoppedChannel()
			if err := d.main(d.stop); err != nil {
				log.Printf("daemon stopped unexpectedly: %s", err)
			}
		}()
		return nil
	}

	return fmt.Errorf("daemon could not be started: %s", "already running")
}

func (d *daemon) Stop() error {
	// TODO: Make thread safe
	if d.Running() {
		d.closeStopChannel()
		d.waitUntilStopped()
		return nil
	}

	return fmt.Errorf("daemon could not be stopped: %s", "not currently running")
}

func (d *daemon) waitUntilStopped() {
	if d.stopped != nil {
		<-d.stopped
	}
}

func (d *daemon) closeStoppedChannel() {
	if d.stopped != nil {
		close(d.stopped) // HINT: Closing makes receivers recieve a zero value
	}
	d.stopped = nil
}

func (d *daemon) closeStopChannel() {
	if d.stop != nil {
		close(d.stop) // HINT: Closing makes receivers recieve a zero value
	}
	d.stop = nil
}
