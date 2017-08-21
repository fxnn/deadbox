package config

const DefaultMaxWorkerTimeoutInSeconds = 100
const DefaultMaxRequestTimeoutInSeconds = 24 * 60 * 60

// Drop configuration, created once per configured drop
type Drop struct {
	// Name identifies this drop uniquely
	Name string

	// ListenAddress defines the the local network address this drop
	// shall listen on.
	ListenAddress string

	// MaxWorkerTimeoutInSeconds limits the time period a worker registration may be considered active without having
	// received an update from the worker.
	MaxWorkerTimeoutInSeconds int

	// MaxRequestTimeoutInSeconds limits the time period during which a request and its response may be processed and
	// retrieved.
	MaxRequestTimeoutInSeconds int
}
