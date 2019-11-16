package config

import "net/url"

func Dummy() *Application {
	dropUrl, _ := url.Parse("https://localhost:" + DefaultPort)
	w := Worker{
		Name:                                "Default Worker",
		DropUrl:                             dropUrl,
		UpdateRegistrationIntervalInSeconds: DefaultUpdateRegistrationIntervalInSeconds,
		RegistrationTimeoutInSeconds:        DefaultRegistrationTimeoutInSeconds,
		PrivateKeySize:                      DefaultPrivateKeySize,
	}
	d := Drop{
		Name:                      "Default Drop",
		ListenAddress:             ":" + DefaultPort,
		MaxWorkerTimeoutInSeconds: DefaultMaxWorkerTimeoutInSeconds,
		PrivateKeySize:            DefaultPrivateKeySize,
		CertificateHosts:          DefaultCertificateHosts,
		CertificateValidFor:       DefaultCertificateValidFor,
	}
	app := &Application{
		DbPath:                             "./",
		PrivateKeyPath:                     "./",
		CertPath:                           "./",
		Workers:                            []Worker{w},
		Drops:                              []Drop{d},
		PublicKeyFingerprintChallengeLevel: DefaultPublicKeyFingerprintChallengeLevel,
		PublicKeyFingerprintLength:         DefaultPublicKeyFingerprintLength,
	}
	return app
}
