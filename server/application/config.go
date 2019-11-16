package application

import (
	"github.com/fxnn/deadbox/server/config"
	"github.com/fxnn/deadbox/server/crypto"
)

func ReadConfig() *config.Application {
	// @todo #4 replace dummy config with configuration mechanism
	cfg := config.Dummy()

	// @todo #31 remove fingerprint-shortcut when dummy config is removed
	_, privateKey, err := ReadOrCreatePrivateKey(cfg.PrivateKeyPath, cfg.Drops[0].Name, cfg.Drops[0].PrivateKeySize, cfg.PublicKeyFingerprintLength,
		cfg.PublicKeyFingerprintChallengeLevel)
	if err != nil {
		panic(err)
	}

	fingerprint, err := crypto.FingerprintPublicKey(&privateKey.PublicKey, cfg.PublicKeyFingerprintChallengeLevel, cfg.PublicKeyFingerprintLength)
	if err != nil {
		panic(err)
	}

	cfg.Workers[0].DropFingerprint = fingerprint

	return cfg
}
