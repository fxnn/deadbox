package config

import (
	"net/url"
)

const (
	DefaultUpdateRegistrationIntervalInSeconds = 10
	DefaultRegistrationTimeoutInSeconds        = 10 * DefaultUpdateRegistrationIntervalInSeconds
)

// Worker configuration, created once per configured worker.
type Worker struct {
	// Name identifies this worker for human users
	Name string

	// DropUrl identifies the drop instances this worker should connect to
	DropUrl *url.URL

	// UpdateRegistrationIntervalInSeconds specifies how often the worker sends a registration update to the drop.
	UpdateRegistrationIntervalInSeconds int

	// RegistrationTimeoutInSeconds specifies how long the registration at the drop is requested to be valid without
	// sending an update.
	RegistrationTimeoutInSeconds int

	// PrivateKeySize is the size of the private RSA key in bits, mostly 2048 oder 4096.
	PrivateKeySize int

	// DropFingerprint allows to validate the Drops TLS certificate using the fingerprint mechanism. If empty, the
	// usual certificate validation is used.
	DropFingerprint string
}
