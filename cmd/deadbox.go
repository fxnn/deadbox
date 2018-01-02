package main

import (
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
		daemons = append(daemons, serveDrop(dp, cfg))
	}

	for _, wk := range cfg.Workers {
		daemons = append(daemons, runWorker(wk, cfg))
	}

	return daemons
}
func waitForShutdownRequest() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

func runWorker(wcfg config.Worker, acfg *config.Application) daemon.Daemon {
	var k = readOrCreatePrivateKeyFile(wcfg.PrivateKeyFile)
	var id = generateWorkerId(k, wcfg.PublicKeyFingerprintLength, wcfg.PublicKeyFingerprintChallengeLevel)

	var b = openDb(acfg, id)
	var d daemon.Daemon = worker.New(wcfg, id, b, k)
	d.OnStop(b.Close)
	d.Start()

	return d
}
func generateWorkerId(privateKeyBytes []byte, fingerprintLength int, challengeLevel uint) string {
	if privateKey, err := crypto.UnmarshalPrivateKeyFromPEMBytes(privateKeyBytes); err != nil {
		panic(fmt.Errorf("couldn't read private key file: %s", err))
	} else if fingerprint, err := crypto.FingerprintPublicKey(&privateKey.PublicKey, challengeLevel, fingerprintLength); err != nil {
		panic(err)
	} else {
		return fingerprint
	}
}

func readOrCreatePrivateKeyFile(fileName string) []byte {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Errorf("couldn't read file %s: %s", fileName, err))
		}

		bytes, err = worker.GeneratePrivateKeyBytes()
		if err != nil {
			panic(fmt.Errorf("couldn't generate private key: %s", err))
		}

		err = ioutil.WriteFile(fileName, bytes, filePermOnlyUserCanReadOrWrite)
		if err != nil {
			panic(fmt.Errorf("couldn't write generated private key to file %s: %s", fileName, err))
		}
	}

	return bytes
}

func serveDrop(dcfg config.Drop, acfg *config.Application) daemon.Daemon {
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
