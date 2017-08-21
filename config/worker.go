package config

import (
	"net/url"
)

const DefaultUpdateRegistrationIntervalInSeconds = 10
const DefaultRegistrationTimeoutInSeconds = 10 * DefaultUpdateRegistrationIntervalInSeconds

// Worker configuration, created once per configured worker.
type Worker struct {
	// Name identifies this worker uniquely
	Name string

	// DropUrl identifies the drop instances this worker should
	// connect to
	DropUrl *url.URL

	// UpdateRegistrationIntervalInSeconds specifies how often the worker sends a registration update to the drop.
	UpdateRegistrationIntervalInSeconds int

	// RegistrationTimeoutInSeconds specifies how long the registration at the drop is requested to be valid without
	// sending an update.
	RegistrationTimeoutInSeconds int
}
