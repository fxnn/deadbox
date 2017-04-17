package config

import "net/url"

// Worker configuration, created once per configured worker.
type Worker struct {
	// Name identifies this worker uniquely
	Name string

	// DropUrl identifies the drop instances this worker should
	// connect to
	DropUrl *url.URL
}
