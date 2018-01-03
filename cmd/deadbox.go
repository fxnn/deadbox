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
	"github.com/fxnn/deadbox/rest"
	"github.com/fxnn/deadbox/worker"
)

const (
	filePermOnlyUserCanReadOrWrite = 0600
	dbFileExtension                = "boltdb"
	privateKeyFileExtension        = "pem"
	certFileExtension              = "pem"
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
	var _, key, err = readOrCreatePrivateKey(acfg.PrivateKeyPath, wcfg.Name, wcfg.PrivateKeySize)
	if err != nil {
		panic(err)
	}
	var id = generateWorkerId(key, wcfg.PublicKeyFingerprintLength, wcfg.PublicKeyFingerprintChallengeLevel)
	var db = openDb(acfg, wcfg.Name)

	var d daemon.Daemon = worker.New(wcfg, id, db, key)
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

func readOrCreatePrivateKey(privateKeyPath string, name string, privateKeySize int) (fileName string, privateKey *rsa.PrivateKey, err error) {
	var bytes []byte

	fileName = privateKeyFileName(privateKeyPath, name)
	bytes, err = ioutil.ReadFile(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			err = fmt.Errorf("couldn't read file %s: %s", fileName, err)
			return
		}

		log.Printf("worker '%s' has no private key, generating one", name)
		bytes, err = worker.GeneratePrivateKeyBytes(privateKeySize)
		if err != nil {
			err = fmt.Errorf("couldn't generate private key: %s", err)
			return
		}

		err = ioutil.WriteFile(fileName, bytes, filePermOnlyUserCanReadOrWrite)
		if err != nil {
			err = fmt.Errorf("couldn't write generated private key to file %s: %s", fileName, err)
			return
		}
	}

	if privateKey, err = crypto.UnmarshalPrivateKeyFromPEMBytes(bytes); err != nil {
		err = fmt.Errorf("couldn't read private key from file %s: %s", fileName, err)
		return
	}
	if privateKey.N.BitLen() != privateKeySize {
		log.Printf("worker '%s' has configured key size '%d', but existing key has size '%d'",
			name, privateKeySize, privateKey.N.BitLen())
	}

	return
}

func serveDrop(dcfg *config.Drop, acfg *config.Application) (daemon daemon.Daemon) {
	db := openDb(acfg, dcfg.Name)
	tls, err := getOrCreateTLSCertFiles(acfg.CertPath, acfg.PrivateKeyPath, dcfg.Name, dcfg.PrivateKeySize, dcfg.CertificateValidFor, dcfg.CertificateHosts)
	if err != nil {
		panic(err)
	}

	daemon = drop.New(dcfg, db, tls)
	daemon.OnStop(db.Close)
	daemon.Start()

	return
}

func getOrCreateTLSCertFiles(certPath string, privateKeyPath string, name string, privateKeySize int, certificateValidFor time.Duration, hosts []string) (rest.TLS, error) {
	privateKeyFile, privateKey, err := readOrCreatePrivateKey(privateKeyPath, name, privateKeySize)
	if err != nil {
		return nil, err
	}

	certFile := certFileName(certPath, name)
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Printf("drop '%s' has no certificate, generating one", name)
		bytes, err := crypto.GenerateCertificateBytes(privateKey, certificateValidFor, hosts)
		if err != nil {
			err = fmt.Errorf("couldn't generate certificate: %s", err)
			return nil, err
		}

		err = ioutil.WriteFile(certFile, bytes, filePermOnlyUserCanReadOrWrite)
		if err != nil {
			err = fmt.Errorf("couldn't write generated certificate to file %s: %s", certFile, err)
			return nil, err
		}
	}

	return rest.NewFileBasedTLS(privateKeyFile, certFile), nil
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

func privateKeyFileName(path string, name string) string {
	return filepath.Join(path, name+".private."+privateKeyFileExtension)
}

func certFileName(path string, dropName string) string {
	return filepath.Join(path, dropName+"."+certFileExtension)
}
