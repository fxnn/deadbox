package daemon

import (
	"fmt"
	"log"
)

type Daemon interface {
	Running() bool
	Start() error
	Stop() error
	OnStop(func() error)
}

type daemon struct {
	main    func(stop <-chan struct{}) error
	stop    chan struct{}
	stopped chan struct{}
	onStop  []func() error
}

func New(main func(stop <-chan struct{}) error) Daemon {
	return &daemon{main, nil, nil, make([]func() error, 0)}
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
			defer d.sendStoppedEvent()
			defer d.closeStopChannel() // HINT: Just in case we terminate abnormally
			defer d.invokeOnStopHandlers()
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
		d.sendStopEvent()
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

func (d *daemon) sendStoppedEvent() {
	d.closeStoppedChannel() // HINT: Closing makes receivers recieve a zero value
}

func (d *daemon) closeStoppedChannel() {
	if d.stopped != nil {
		close(d.stopped)
	}
	d.stopped = nil
}

func (d *daemon) sendStopEvent() {
	d.closeStopChannel() // HINT: Closing makes receivers recieve a zero value
}

func (d *daemon) closeStopChannel() {
	if d.stop != nil {
		close(d.stop)
	}
	d.stop = nil
}

func (d *daemon) OnStop(handler func() error) {
	d.onStop = append(d.onStop, handler)
}

func (d *daemon) invokeOnStopHandlers() {
	for _, l := range d.onStop {
		if err := l(); err != nil {
			log.Printf("daemon's onStop handler failed: %s", err)
		}
	}
}
