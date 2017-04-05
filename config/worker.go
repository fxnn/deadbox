package config

import "net/url"

// Worker configuration, created once per configured worker.
type Worker struct {
	// Name identifies this worker uniquely
	Name string

	// DropUrls identifies the drop instances this worker should
	// connect to
	DropUrls []*url.URL
}
