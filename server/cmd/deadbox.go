package main

import (
	"crypto/rsa"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fxnn/deadbox/server/application"
	"github.com/fxnn/deadbox/server/config"
	"github.com/fxnn/deadbox/server/crypto"
	"github.com/fxnn/deadbox/server/daemon"
	"github.com/fxnn/deadbox/server/drop"
	"github.com/fxnn/deadbox/server/worker"
)

func main() {
	var cfg = application.ReadConfig()
	daemons := startDaemons(cfg)

	waitForShutdownRequest()

	log.Println("shutting down gracefully")
	shutdownDaemons(daemons)
}

func shutdownDaemons(daemons []daemon.Daemon) {
	for _, d := range daemons {
		if err := d.Stop(); err != nil {
			log.Println(err)
		}
	}
}

func startDaemons(cfg *config.Application) []daemon.Daemon {
	var daemons = make([]daemon.Daemon, 0, len(cfg.Drops)+len(cfg.Workers))

	for _, dp := range cfg.Drops {
		daemons = append(daemons, serveDrop(&dp, cfg))
	}

	for _, wk := range cfg.Workers {
		daemons = append(daemons, runWorker(&wk, cfg))
	}

	return daemons
}

func waitForShutdownRequest() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

func runWorker(wcfg *config.Worker, acfg *config.Application) daemon.Daemon {
	var _, key, err = application.ReadOrCreatePrivateKey(acfg.PrivateKeyPath, wcfg.Name, wcfg.PrivateKeySize, acfg.PublicKeyFingerprintLength, acfg.PublicKeyFingerprintChallengeLevel)
	if err != nil {
		panic(err)
	}
	var id = generateWorkerId(key, acfg.PublicKeyFingerprintLength, acfg.PublicKeyFingerprintChallengeLevel)
	var db = application.OpenDb(acfg, wcfg.Name)

	var d daemon.Daemon = worker.New(wcfg, id, db, key, acfg.PublicKeyFingerprintLength, acfg.PublicKeyFingerprintChallengeLevel)
	d.OnStop(db.Close)
	d.Start()

	return d
}

func generateWorkerId(privateKey *rsa.PrivateKey, fingerprintLength uint, challengeLevel uint) string {
	if fingerprint, err := crypto.FingerprintPublicKey(&privateKey.PublicKey, challengeLevel, fingerprintLength); err != nil {
		panic(err)
	} else {
		return fingerprint
	}
}

func serveDrop(dcfg *config.Drop, acfg *config.Application) (daemon daemon.Daemon) {
	db := application.OpenDb(acfg, dcfg.Name)
	tls, err := application.ReadOrCreateTLSCertFiles(
		acfg.CertPath,
		acfg.PrivateKeyPath,
		dcfg.Name,
		dcfg.PrivateKeySize,
		dcfg.CertificateValidFor,
		dcfg.CertificateHosts,
		acfg.PublicKeyFingerprintLength,
		acfg.PublicKeyFingerprintChallengeLevel)
	if err != nil {
		panic(err)
	}

	daemon = drop.New(dcfg, db, tls)
	daemon.OnStop(db.Close)
	daemon.Start()

	return
}
