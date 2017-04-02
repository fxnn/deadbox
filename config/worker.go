package config

import "net/url"

// Worker configuration, created once per configured worker.
type Worker struct {

	// Name identifies this agent uniquely
	Name string

	// DropUrls identifies the drop instances this agent should
	// connect to
	DropUrls []*url.URL
}
