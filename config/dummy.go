package config

import "net/url"

func Dummy() *Application {
	dropUrl, _ := url.Parse("http://localhost:" + DefaultPort)
	w := Worker{
		Name:    "Default Worker",
		DropUrl: dropUrl,
		UpdateRegistrationIntervalInSeconds: DefaultUpdateRegistrationIntervalInSeconds,
		RegistrationTimeoutInSeconds:        DefaultRegistrationTimeoutInSeconds,
		PublicKeyFingerprintChallengeLevel:  DefaultPublicKeyFingerprintChallengeLevel,
		PublicKeyFingerprintLength:          DefaultPublicKeyFingerprintLength,
	}
	d := Drop{
		Name:                      "Default Drop",
		ListenAddress:             ":" + DefaultPort,
		MaxWorkerTimeoutInSeconds: DefaultMaxWorkerTimeoutInSeconds,
	}
	app := &Application{
		DbPath:         "./",
		PrivateKeyPath: "./",
		Workers:        []Worker{w},
		Drops:          []Drop{d},
	}
	return app
}
