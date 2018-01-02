package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fxnn/deadbox/config"
	"github.com/fxnn/deadbox/crypto"
	"github.com/fxnn/deadbox/daemon"
	"github.com/fxnn/deadbox/drop"
	"github.com/fxnn/deadbox/worker"
)

const (
	filePermOnlyUserCanReadOrWrite = 0600
	dbFileExtension                = "boltdb"
	privateKeyFileExtension        = "pem"
)

func main() {

	// @todo #4 replace dummy config with configuration mechanism
	var cfg = config.Dummy()
	daemons := startDaemons(cfg)

	waitForShutdownRequest()

	log.Println("Shutting down gracefully")
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
	var k = readOrCreatePrivateKey(acfg, wcfg)
	var id = generateWorkerId(k, wcfg.PublicKeyFingerprintLength, wcfg.PublicKeyFingerprintChallengeLevel)

	var b = openDb(acfg, wcfg.Name)
	var d daemon.Daemon = worker.New(wcfg, id, b, k)
	d.OnStop(b.Close)
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

func readOrCreatePrivateKey(acfg *config.Application, wcfg *config.Worker) *rsa.PrivateKey {
	fileName := privateKeyFileName(acfg.PrivateKeyPath, wcfg.Name)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Errorf("couldn't read file %s: %s", fileName, err))
		}

		log.Printf("worker '%s' has no private key, generating one", wcfg.Name)
		bytes, err = worker.GeneratePrivateKeyBytes(wcfg.PrivateKeySize)
		if err != nil {
			panic(fmt.Errorf("couldn't generate private key: %s", err))
		}

		err = ioutil.WriteFile(fileName, bytes, filePermOnlyUserCanReadOrWrite)
		if err != nil {
			panic(fmt.Errorf("couldn't write generated private key to file %s: %s", fileName, err))
		}
	}

	if privateKey, err := crypto.UnmarshalPrivateKeyFromPEMBytes(bytes); err != nil {
		panic(fmt.Errorf("couldn't read private key from file %s: %s", fileName, err))
	} else {
		if privateKey.N.BitLen() != wcfg.PrivateKeySize {
			log.Printf("worker '%s' has configured key size '%d', but existing key has size '%d'",
				wcfg.Name, wcfg.PrivateKeySize, privateKey.N.BitLen())
		}

		return privateKey
	}
}

func serveDrop(dcfg *config.Drop, acfg *config.Application) daemon.Daemon {
	var b = openDb(acfg, dcfg.Name)
	var d daemon.Daemon = drop.New(dcfg, b)
	d.OnStop(b.Close)
	d.Start()

	return d
}

func openDb(cfg *config.Application, name string) *bolt.DB {
	boltOptions := &bolt.Options{Timeout: 10 * time.Second}

	fileName := dbFileName(cfg, name)
	db, err := bolt.Open(fileName, 0660, boltOptions)
	if err != nil {
		panic(fmt.Errorf(
			"couldn't open bolt DB %s: %s",
			fileName, err,
		))
	}
	log.Println("Database opened:", fileName)

	return db
}

func dbFileName(cfg *config.Application, name string) string {
	return filepath.Join(cfg.DbPath, name+"."+dbFileExtension)
}

func privateKeyFileName(path string, workerName string) string {
	return filepath.Join(path, workerName+"."+privateKeyFileExtension)
}
