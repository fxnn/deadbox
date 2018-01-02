package config

import (
	"net/url"
)

const DefaultUpdateRegistrationIntervalInSeconds = 10
const DefaultRegistrationTimeoutInSeconds = 10 * DefaultUpdateRegistrationIntervalInSeconds

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

	// PrivateKeyFile refers to a file containing the ASN.1 PKCS#1 DER encoded private RSA key.
	// If not existant, it will be created.
	PrivateKeyFile string

	// PublicKeyFingerprintLength influences the length of the public keys fingerprint. The greater the length, the
	// more reliable the fingerprint is, but the harder it is to remember for human users.
	PublicKeyFingerprintLength uint

	// PublicKeyFingerprintChallengeLevel influences the time it takes to generate the fingerprint. The greater the
	// level, the more secure the fingerprint is against pre-image attacks, but the longer it takes to generate and
	// validate the fingerprint.
	PublicKeyFingerprintChallengeLevel uint
}
