package daemon

import "fmt"

type Daemon interface {
	Running() bool
	Start() error
	Stop() error
}

type daemon struct {
	main func(stop <-chan struct{}) error
	stop chan struct{}
}

func New(main func(stop <-chan struct{}) error) Daemon {
	return &daemon{main, nil}
}

func (d *daemon) Running() bool {
	return d.stop != nil
}

func (d *daemon) Start() error {
	// TODO: Make thread safe
	if !d.Running() {
		if d.main == nil {
			return fmt.Errorf("no main set")
		}

		d.stop = make(chan struct{})
		go func() {
			defer d.closeStopChannel()
			if err := d.main(d.stop); err != nil {
				fmt.Println(err)
			}
		}()
		return nil
	}

	return fmt.Errorf("already running")
}

func (d *daemon) Stop() error {
	// TODO: Make thread safe
	if d.Running() {
		d.closeStopChannel()
		return nil
	}

	return fmt.Errorf("not currently running")
}

func (d *daemon) closeStopChannel() {
	close(d.stop) // HINT: Closing makes receivers recieve a zero value
	d.stop = nil
}
